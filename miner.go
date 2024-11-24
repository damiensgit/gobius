package main

import (
	"container/list"
	"context"
	"database/sql"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gobius/arbius/engine"
	"gobius/client"
	task "gobius/common"
	"gobius/config"
	"gobius/ipfs"
	"gobius/models"
	"gobius/tui"
	"gobius/utils"
	"io"
	"log"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

//go:embed sql/sqlite/migrations/*.sql
var embedMigrations embed.FS

// Generate the required go file from the engine contract
// This requires the official arbius project to be cloned and available TODO: make this better?
////go:generate ./bin/solc-latest --base-path '..' --include-path '..\arbiusv3\contract\node_modules' --bin --abi ./contracts/V2_EngineV3.sol -o build --overwrite --evm-version london
//go:generate ./bin/solc --base-path './external/arbius/' --include-path './external/arbius/contract/node_modules' --bin --abi ./contracts/V2_EngineV4.sol -o build --overwrite --evm-version london

//go:generate abigen --bin=./build/V2_EngineV3.bin --pkg enginev3 --abi ./build/V2_EngineV3.abi --out ./arbius/engine/enginev3.go
//go:generate abigen --bin=./build/IBaseToken.bin --pkg basetoken --abi ./build/IBaseToken.abi --out ./arbius/basetoken/basetoken.go

//./bin/solc-latest --base-path './contracts' --bin --abi .\contracts\BulkTasksNova.sol -o build --overwrite --evm-version london
//go:generate ./bin/solc-latest --include-path './contracts' --bin --abi ./contracts/DelegatedValidatorV2.sol -o build --overwrite --evm-version london
//go:generate abigen --bin=./build/DelegatedValidatorV2.bin --pkg delegatedvalidator --abi ./build/DelegatedValidatorV2.abi --out ./arbius/delegatedminer/delegatedminer.go
//abigen --bin=./build/BulkTasks.bin --pkg bulktasks --abi ./build/BulkTasks.abi --out ./arbius/bulktasks/bulktasks.go

const WETH_ADDRESS = "0x722e8bdd2ce80a4422e880164f2079488e115365"

type TaskSubmitted struct {
	TaskId task.TaskId
	TxHash common.Hash
}

const (
	maxTasks        = 50 // Maximum number of tasks to store
	maxExitAttempts = 3  // Maximum number of attempts to exit the application
)

func initLogging(file string, level zerolog.Level, consoleWriter io.Writer) (logFile io.WriteCloser, logger zerolog.Logger) {
	// logFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatalf("Failed opening log file: %s", err)
	// }

	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.TimeFieldFormat = time.RFC3339Nano

	fileLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     5, // days
	}

	if consoleWriter == nil {
		consoleWriter = os.Stderr
	}

	multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: consoleWriter, TimeFormat: "15:04:05.000000000"}, fileLogger)
	//multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.000000000"}, fileLogger)

	logger = zerolog.New(multi).Level(level).With().Timestamp().Logger()

	return fileLogger, logger
}

func setupCloseHandler(logger *zerolog.Logger) (ctx context.Context) {
	// Create a cancellation context.
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		exitCounter := maxExitAttempts

		for {
			<-c

			switch {
			case exitCounter < 3 && exitCounter > 0:
				logger.Error().Int("attempt", maxExitAttempts+1-exitCounter).Msgf("Already shutting down. Will force exit after %d more attempt(s).", exitCounter)
				exitCounter--
			case exitCounter == 0:
				logger.Error().Msgf("Ctrl+C pressed %d times. Forcing exit.", maxExitAttempts)
				os.Exit(0)
			default:
				logger.Info().Msg("Ctrl+C detected. Stopping process...")
				cancel()
				exitCounter--
			}
		}
	}()

	return ctx
}

// type TaskManagerI interface {
// 	GetCurrentReward() (*big.Int, error)
// 	GetTotalTasks() int64
// 	GetClaims() int64
// 	GetSolutions() int64
// 	GetCommitments() int64
// }

// TODO: put this into base config ?
var minerVersion = big.NewInt(3)

type Miner struct {
	services  *Services
	ctx       context.Context
	mu        sync.Mutex
	gpura     *utils.RunningAverage
	validator IValidator
	wg        *sync.WaitGroup
	sync.RWMutex
}

type IValidator interface {
	SignalCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) error
	SubmitSolution(validator common.Address, taskId task.TaskId, cid []byte) error
	//ValidatorDeposit(depositAmount *big.Int) (*types.Transaction, error)
	GetNextValidatorAddress() common.Address
	//ValidatorAddress() common.Address
	InitiateValidatorWithdraw(validator common.Address, amount float64) error
	ValidatorWithdraw(validator common.Address) error
	CancelValidatorWithdraw(validator common.Address, count int64) error
	BulkClaim(taskIds [][32]byte) (*types.Receipt, error)
	BatchCommitments() error
	BatchSolutions() error
	VoteOnContestation(validator common.Address, taskId task.TaskId, yeah bool) error
	SubmitContestation(validator common.Address, taskId task.TaskId) error
}

// type TransactionTask struct {
// 	taskId    task.TaskId
// 	taskType  int // 0 = commitment, 1 = solution
// 	cid       []byte
// 	commiment [32]byte
// }

// connects miner account to the engine contract
// you can only have one miner account to one enginev2 for submitting tasks and solutions
// you can have any account submit claims so that will be done elsewhere
func NewMinerEngine(ctx context.Context, validator IValidator, wg *sync.WaitGroup) (*Miner, error) {

	ra := utils.NewRunningAverage(15 * time.Minute)

	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	miner := &Miner{
		services: services,
		ctx:      ctx,
		gpura:    ra,
		mu:       sync.Mutex{},
		wg:       wg,
	}

	miner.validator = validator

	return miner, nil
}

// Define the function parameters
type SubmitTaskParams struct {
	Version uint8
	Owner   common.Address
	Model   [32]byte
	Fee     *big.Int
	Input   []byte
}

func (m *Miner) AvergeGPUTime() time.Duration {
	return m.gpura.Average()
}

func (m *Miner) SolveTask(taskId task.TaskId, tx *types.Transaction, gpu *task.GPU, validateOnly bool) ([]byte, error) {
	// authCopy := new(bind.TransactOpts)
	// *authCopy = *m.services.OwnerAccount.Auth
	// // Assign the copied context to the copied auth object
	// authCopy.Context = context.Background()

	//startOfSolve := time.Now()
	taskIdStr := taskId.String()

	//m.services.Logger.Error().Uint64("nonce", m.validator.Nonce()).Msg("START OF SOLUTION SOLVE")

	// from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	// if err != nil {
	// 	from, _ = types.Sender(types.HomesteadSigner{}, tx)
	// }

	/*from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return nil, err
	}

	var taskData *SubmitTaskParams

	if from.String() != m.services.SenderOwnerAccount.Address.String() {
		m.services.Logger.Debug().Str("sender", from.String()).Msg("using decoder to get submit task data")

		if m.services.Config.WhitelistTasks {
			if _, ok := whitelist[*tx.To()]; ok {
				m.services.Logger.Debug().Str("taskid", taskIdStr).Str("txhash", tx.Hash().String()).Msgf("whitelisted task - skipping (to: %s)", tx.To().String())
				return nil, nil
			}
		}

		taskData, err = m.DecodeSubmitTask(tx, taskId)
		if err != nil {
			m.services.Logger.Error().Err(err).Str("taskid", taskIdStr).Str("txhash", tx.Hash().String()).Msgf("could not decode transaction: %s", common.Bytes2Hex(tx.Data()))
			return nil, err
		}
		if !validateOnly && m.services.Config.StealTasks {
			if taskData.Owner.String() != m.validator.ValidatorAddress().String() {
				if _, ok := m.services.Config.Strategies.Snipe.Targets[taskData.Owner]; !ok {
					m.services.Logger.Debug().Str("taskid", taskIdStr).Msgf("task owner %s not in sniper target list, skipping", taskData.Owner.String())
					return nil, nil
				}
			}
		}
	} else {
		taskData = m.services.AutoMineParams
	}*/

	taskData := m.services.AutoMineParams

	inputRaw := string(taskData.Input)

	//var input models.Inner
	var result map[string]interface{}
	err := json.Unmarshal(taskData.Input, &result)
	//err = json.Unmarshal(taskData.Input, &input)
	if err != nil {
		data := string(taskData.Input)

		// Remove the leading 0x
		hexString := strings.TrimPrefix(data, "0x")

		decoded, err := hex.DecodeString(hexString)
		if err != nil {
			hexStr := hex.EncodeToString(taskData.Input)
			m.services.Logger.Error().Err(err).Str("txhash", tx.Hash().String()).Str("data", hexStr).Msg("could not decode input")
			return nil, err
		}
		err = nil
		var errHex error
		// try decode upto 5 times
		for i := 0; i < 5; i++ {
			err = json.Unmarshal(decoded, &result)
			if err != nil {
				decoded, errHex = hex.DecodeString(string(decoded))
				if errHex != nil {
					return nil, errHex
				}
			} else {
				break
			}
		}

		if err != nil {
			m.services.Logger.Error().Err(err).Str("txhash", tx.Hash().String()).Str("data", string(decoded)).Msg("could not unmarshal input")
			return nil, err
		}
	}

	m.services.Logger.Debug().Str("input", inputRaw).Msg("decoded information")

	taskInfo, err := m.services.Engine.LookupTask(taskId)
	if err != nil {
		m.services.Logger.Error().Err(err).Msg("could not lookup task")
		return nil, err
	}

	modelId := common.Bytes2Hex(taskInfo.Model[:])
	model := models.EnabledModels.FindModel(modelId)
	if model == nil {
		m.services.Logger.Error().Str("model", modelId).Err(err).Msg("could not find model")
		return nil, err
	}

	hydrated, err := model.HydrateInput(result)

	prompt, ok := hydrated.(models.Kandinsky2Prompt)

	if !ok {
		m.services.Logger.Error().Msg("could not hydrate input")
		return nil, err
	}

	// set the seed to the task id
	prompt.Input.Seed = taskId.TaskId2Seed()

	output, err := json.Marshal(prompt)
	if err != nil {
		m.services.Logger.Error().Err(err).Msg("could not marshal output")
		return nil, err
	}

	m.services.Logger.Debug().Str("output", string(output)).Msg("sending task to gpu")

	var cid []byte
	if m.services.Config.EvilMode {
		cid, _ = hex.DecodeString("12206666666666666666666666666666666666666666666666666666666666666666")
		m.services.Logger.Warn().Str("cid", "0x"+hex.EncodeToString(cid)).Msg("evil mode enabled")
		duration := time.Duration(m.services.Config.EvilModeMinTime+rand.Intn(m.services.Config.EvilModeRandInt)) * time.Millisecond
		time.Sleep(duration)
	} else {
		start := time.Now()
		if gpu.Mock {
			data, err := gpu.GetMockCid(taskIdStr, prompt)
			if err != nil {
				return nil, err
			}
			cid, err = ipfs.GetIPFSHashFast(data)
			if err != nil {
				return nil, err
			}
			//	m.services.Logger.Warn().Str("cid", "0x"+hex.EncodeToString(cid)).Msg("mock gpu cid")
		} else {
			cid, err = model.GetCID(gpu, taskIdStr, prompt)
		}
		elapsed := time.Since(start)
		m.gpura.Add(elapsed)
		if err != nil {
			m.services.Logger.Error().Err(err).Msg("error on gpu, incrementing error counter")
			gpu.IncrementErrorCount()
			return nil, err
		}
		m.services.Logger.Debug().Str("cid", "0x"+hex.EncodeToString(cid)).Str("elapsed", elapsed.String()).Str("average", m.gpura.Average().String()).Msg("gpu finished & returned result")
	}

	if validateOnly {
		return cid, nil
	}

	validator := m.validator.GetNextValidatorAddress()

	commitmentFunc := func() error {

		commitment, err := utils.GenerateCommitment(validator, taskId, cid)
		if err != nil {
			m.services.Logger.Error().Err(err).Msg("error generating commitment hash")
			return err
		}

		if m.services.Config.CheckCommitment {
			m.services.Logger.Debug().Str("taskid", taskIdStr).Str("commitment", "0x"+hex.EncodeToString(commitment[:])).Msg("checking for existing task commitment")

			block, err := m.services.Engine.Engine.Commitments(nil, commitment)
			if err != nil {
				m.services.Logger.Error().Err(err).Msg("error getting commitment")
				return err
			}

			blockNo := block.Uint64()
			if blockNo > 0 {
				m.services.Logger.Warn().Str("taskid", taskIdStr).Uint64("block", blockNo).Str("commitment", "0x"+hex.EncodeToString(commitment[:])).Msg("commitment already exists for task")
				return nil
			}
		}

		m.validator.SignalCommitment(validator, taskId, commitment)

		return nil
	}

	err = commitmentFunc()
	if err != nil {
		m.services.Logger.Warn().Msg("commitment failed so not sending solution")
		return nil, err
	}

	go func() {
		m.validator.SubmitSolution(validator, taskId, cid)
		//m.services.Logger.Info().Str("taskid", taskIdStr).Str("totalsolvetime", time.Since(startOfSolve).String()).Msg("solution time")
	}()

	return cid, nil
}

// TODO: move this into miner engine wrapper?
// TODO: or better yet move to to model itself
func validateModel(model models.ModelInterface, gpu *task.GPU, logger *zerolog.Logger) error {

	// model, found := baseCfg.Models["kandinsky2"]
	// if !found {
	// 	panic("model not found - check config/config.json")
	// }

	// m := models.EnabledModels.FindModel(model.ID)
	// if m == nil {
	// 	panic("model not found in enabled models")
	// }

	testPrompt := models.Kandinsky2Prompt{
		Input: models.Inner{
			Prompt: "render a cat in the style of kandinsky",
			Height: 768,
			Width:  768,
			Seed:   1337,
		},
	}

	cid, err := model.GetCID(gpu, "startup-test-taskid", testPrompt)
	if err != nil {
		return err
	}

	expected := "0x12200f8c99111abf301ceb8965af7b111c77bcd6e1903c0c713c4b610665dd270be3"
	cidStr := "0x" + hex.EncodeToString(cid)
	if cidStr == expected {
		logger.Info().Str("model", model.GetID()).Str("cid", cidStr).Str("expected", expected).Msg("model CID matches expected CID")
	} else {
		logger.Error().Str("model", model.GetID()).Str("cid", cidStr).Str("expected", expected).Msg("model CID does not match expected CID")
		return errors.New("model CID does not match expected CID")
	}

	return nil
}

func average(times []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, time := range times {
		total += time
	}
	return total / time.Duration(len(times))
}

func findFastestPrompt(model models.ModelInterface, gpus []*task.GPU, logger *zerolog.Logger) error {

	inputs := []string{
		"1 pixel",
		"1 red pixel center",
		"1010101101010101010011111010101010101011010010111010010110101",
		"0x101111101",
		"black void",
		"all black",
		"mono black",
		"black background black foreground black image",
		"black",
		"white",
		"empty",
		"box",
		"square",
		"red",
		"single color",
		"cat",
		"dog",
		"zero",
		"0",
		"generate nothing",
		"do nothing",
		"null",
	}

	logger.Info().Str("model", model.GetID()).Msg("Validating model")

	testPrompt := models.Kandinsky2Prompt{
		Input: models.Inner{
			Height: 768,
			Width:  768,
		},
	}

	fastestTimeSeen := time.Duration(0)
	fastestPrompt := ""

	for _, input := range inputs {
		var wg sync.WaitGroup
		times := make([]time.Duration, len(gpus))

		testPrompt.Input.Prompt = input

		for i, gpu := range gpus {
			wg.Add(1)

			go func(i int, gpu *task.GPU) {
				defer wg.Done()

				start := time.Now()
				_, err := model.GetCID(gpu, "startup-test-taskid", testPrompt)
				if err != nil {
					logger.Error().Err(err).Msg("Error generating CID")
					return
				}
				timeTaken := time.Since(start)

				times[i] = timeTaken

				logger.Info().Int("gpu", gpu.ID).Str("input", input).Str("duration", timeTaken.String()).Msg("cid generated")
			}(i, gpu)
		}

		wg.Wait()

		averageTime := average(times)

		logger.Info().Str("input", input).Str("average", averageTime.String()).Msg("all gpus completed input")

		if fastestTimeSeen == 0 || averageTime < fastestTimeSeen {
			fastestTimeSeen = averageTime
			fastestPrompt = input
		}
	}

	logger.Info().Str("fastest_prompt", fastestPrompt).Str("duration", fastestTimeSeen.String()).Msg("Fastest prompt found")

	return nil
}

func validateGpus(model models.ModelInterface, gpus []*task.GPU, logger *zerolog.Logger) error {

	//findFastestPrompt(modelToMine, gpus, &logger)

	logger.Info().Str("model", model.GetID()).Msg("Validating model on gpu(s)")

	mu := sync.Mutex{}
	fastestTimeSeen := time.Duration(0)

	var wg sync.WaitGroup
	times := make([]time.Duration, len(gpus))

	for i, gpu := range gpus {
		wg.Add(1)

		go func(i int, gpu *task.GPU) {
			defer wg.Done()

			start := time.Now()

			err := validateModel(model, gpu, logger)

			if err != nil {
				logger.Fatal().Msgf("error validating the model on gpu #%d: %s", gpu.ID, gpu.Url)
			}

			timeTaken := time.Since(start)

			times[i] = timeTaken
			mu.Lock()
			if fastestTimeSeen == 0 || timeTaken < fastestTimeSeen {
				fastestTimeSeen = timeTaken
			}
			mu.Unlock()

			logger.Info().Int("gpu", gpu.ID).Str("duration", timeTaken.String()).Msg("cid generated")
		}(i, gpu)
	}

	wg.Wait()

	averageTime := average(times)

	logger.Info().Str("average", averageTime.String()).Str("fastest", fastestTimeSeen.String()).Msg("all gpus completed validation")

	return nil
}

// TODO: config
// TODO: add init methods
func main() {
	var appQuitWG sync.WaitGroup

	// miner := NewMiner(4000 * time.Millisecond)
	// miner.AddJob(NewTask())
	// miner.isGPUIdle(250 * time.Millisecond)
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	skipValidation := flag.Bool("skipvalidation", false, "Skip safety checks and validation of the model and miner version")
	logLevel := flag.Int("loglevel", 1, "Skip safety checks and validation of the model and miner version")
	testnetType := flag.Int("testnet", 0, "Run using testnet - 1 = local, 2 = nova testnet")
	mockGPUs := flag.Int("mockgpus", 0, "mock gpus for testing")
	taskScanner := flag.Int("taskscanner", 0, "scan blocks for unsolved tasks")

	//cmdtaskId := flag.String("taskid", "", "")
	// cmdtxHash := flag.String("taskhash", "", "")
	flag.Parse()

	// Create a set of the flags that were set
	setFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		setFlags[f.Name] = true
	})

	cfg, err := config.InitAppConfig(*configPath, *testnetType)
	if err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	// Check if logLevel was set
	if !setFlags["loglevel"] {
		// If logLevel was not set on the command line, use the value from the config
		*logLevel = int(cfg.LogLevel)
	}

	logWriter := &tui.LogViewWriter{
		App: nil,
	}
	//	_ = writer

	logFile, logger := initLogging(cfg.LogPath, zerolog.Level(*logLevel), logWriter)

	defer logFile.Close()

	if cfg.EvilMode {
		logger.Warn().Msg("TESTNET MODE ENABLED - DO NOT USE ON MAINNET")
	}

	appQuit := setupCloseHandler(&logger)

	rpcClient, err := client.NewClient(cfg.Blockchain.RPCURL, appQuit, cfg.Blockchain.EthersGas, cfg.Blockchain.BasefeeX, cfg.Blockchain.ForceGas, cfg.Blockchain.GasOverride)

	if err != nil {
		logger.Error().Err(err).Msgf("error connecting to RPC: %s", cfg.Blockchain.RPCURL)
		return
	}

	txRpcClient := rpcClient
	if cfg.Blockchain.SenderRPCURL != "" {
		txRpcClient, err = client.NewClient(cfg.Blockchain.SenderRPCURL, appQuit, cfg.Blockchain.EthersGas, cfg.Blockchain.BasefeeX, cfg.Blockchain.ForceGas, cfg.Blockchain.GasOverride)

		if err != nil {
			logger.Error().Err(err).Msgf("error connecting to RPC: %s", cfg.Blockchain.SenderRPCURL)
			return
		}
	}

	var clients []*client.Client
	for _, curl := range cfg.Blockchain.ClientRPCURLs {
		c, err := client.NewClient(curl, appQuit, cfg.Blockchain.EthersGas, cfg.Blockchain.BasefeeX, cfg.Blockchain.ForceGas, cfg.Blockchain.GasOverride)

		if err != nil {
			logger.Error().Err(err).Msgf("error connecting to RPC: %s", curl)
			return
		}

		clients = append(clients, c)
	}

	logger.Info().Str("database", cfg.DBPath).Msg("using database")

	sqlite, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		panic(err)
	}
	_, err = sqlite.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatal(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	if err := goose.Up(sqlite, "sql/sqlite/migrations"); err != nil {
		logger.Fatal().Err(err).Msg("database migration error")
	}

	// _, err = sqlite.Exec("VACUUM")
	// if err != nil {
	// 	logger.Fatal().Err(err).Msg("error on database VACUUM")
	// }

	appContext, err := NewApplicationContext(rpcClient, txRpcClient, clients, sqlite, &logger, cfg, context.Background())

	if err != nil {
		logger.Fatal().Err(err).Msg("could not create application context")
	}

	manager, err := NewBatchTransactionManager(appContext, &appQuitWG)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not create transaction manager")
	}

	// Get non-flag arguments
	// TODO: Make more go idiomatic. Move this to cmd/gobius/* structure
	args := flag.Args()
	if len(args) > 0 {
		command := args[0]

		switch command {
		case "exportconfig":
			if len(args) > 1 {
				filename := args[1]
				err := cfg.ExportConfig(filename)
				if err != nil {
					logger.Fatal().Err(err).Msgf("error exporting config to %s", filename)
				}
				logger.Info().Msgf("config exported to %s", filename)
			} else {
				log.Fatalf("exportconfig requires filename to export to parameter")
			}
		case "unsolvedimport":
			if len(args) > 1 {
				filename := args[1]
				importUnsolvedTasks(filename, &logger, appContext)
			} else {
				log.Fatalf("unsolvedimport requires filename to import parameter")
			}
		case "unclaimedimport":
			if len(args) > 1 {
				filename := args[1]
				importUnclaimedTasks(filename, &logger, appContext)
			} else {
				log.Fatalf("unsolvedimport requires filename to import parameter")
			}
		case "verifyclaims":
			verifyClaims(&logger, appContext)

		case "taskcheck":
			taskCheck(&logger, appContext)

		case "verifysolutions":
			verifySolutions(appContext)

		case "blockmonitor":
			blockMonitor(appContext, rpcClient)

		case "cleantaskdata":
			cleanTaskData(appContext)

		case "claimbatchinfo":
			getBatchPricingInfo(appContext)

		case "fundtaskwallets":
			var amount, minbal float64
			if len(args) == 3 {
				amount, err = strconv.ParseFloat(args[1], 64)
				if err != nil {
					log.Fatalf("Invalid amount: %v", err)
				}
				minbal, err = strconv.ParseFloat(args[2], 64)
				if err != nil {
					log.Fatalf("Invalid minbal: %v", err)
				}
			} else {
				log.Fatalf("fundtaskwallets requires amount to send to each wallet and min bal")
			}
			fundTaskWallets(appContext, amount, minbal)

		case "initiatewithdrawall", "completewithdrawall":
			var amount float64
			var validator common.Address

			if len(args) == 3 {
				amount, err = strconv.ParseFloat(args[1], 64)
				if err != nil {
					logger.Fatal().Msgf("Invalid amount: %v", err)
				}
				if !common.IsHexAddress(args[2]) {
					logger.Fatal().Msgf("Invalid validator address: %s", args[2])
				}
				validator = common.HexToAddress(args[2])
			} else {
				logger.Fatal().Msgf("initiatewithdrawall requires amount (can be set to 0 to withdraw full amount) and validator address")
			}

			switch command {

			case "initiatewithdrawall":
				manager.InitiateValidatorWithdraw(validator, amount)
			case "completewithdrawall":
				manager.ValidatorWithdraw(validator)
			}
		case "totalstaked":
			var validator common.Address
			validator = common.HexToAddress(args[1])

			manager.TotalStaked(validator)

		case "cancelwithdraw":
			var index int64
			var validator common.Address

			if len(args) == 3 {
				index, err = strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					logger.Fatal().Msgf("Invalid index: %v", err)
				}
				if !common.IsHexAddress(args[2]) {
					logger.Fatal().Msgf("Invalid validator address: %s", args[2])
				}
				validator = common.HexToAddress(args[2])
			} else {
				logger.Fatal().Msgf("cancelwithdraw requires index and validator address")
			}
			manager.CancelValidatorWithdraw(validator, index)
		case "voteoncontestation":
			var taskid task.TaskId
			var yea bool
			var validator common.Address

			if len(args) == 4 {
				taskid, err = task.ConvertTaskIdString2Bytes(args[1])
				if err != nil {
					log.Fatalf("Invalid taskid: %v", err)
				}
				yea, err = strconv.ParseBool(args[2])
				if err != nil {
					log.Fatalf("Invalid boolean value: %v", err)
				}
				if !common.IsHexAddress(args[3]) {
					logger.Fatal().Msgf("Invalid validator address: %s", args[3])
				}
				validator = common.HexToAddress(args[3])
			} else {
				log.Fatalf("voteoncontestation requires taskid and true/false")
			}

			manager.VoteOnContestation(validator, taskid, yea)

		case "depositmonitor":
			var endBlock, startBlock int64 = 0, -1
			if len(args) == 2 {
				startBlock, err = strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					log.Fatalf("Invalid start block: %v", err)
				}
			} else if len(args) == 3 {
				startBlock, err = strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					log.Fatalf("Invalid start block: %v", err)
				}
				endBlock, err = strconv.ParseInt(args[2], 10, 64)
				if err != nil {
					log.Fatalf("Invalid end block: %v", err)
				}
			} else {
				log.Fatalf("depositmonitor requires startblock and optional endblock")
			}

			depositMonitor(appContext, rpcClient, startBlock, endBlock)

		default:
			log.Fatalf("unknown command: %s", command)

		}
		return
	}

	// Use the mock IPFS client for now
	// TODO: change this based on strategy in config
	ipfsClient, err := ipfs.NewMockIPFSClient(*cfg, true)

	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to IPFS")
	}

	models.InitEnabledModels(ipfsClient, cfg)

	modelToMine := models.EnabledModels.FindModel(cfg.Strategies.Model)
	if modelToMine == nil {
		logger.Error().Str("model", cfg.Strategies.Model).Msg("model specified in config was not found in enabled models")
		return
	}

	// ui := tui.InitialModel(manager)

	// go func() {
	// 	p := tea.NewProgram(ui, tea.WithAltScreen())
	// 	logWriter.App = p
	// 	if _, err := p.Run(); err != nil {
	// 		fmt.Printf("Alas, there's been an error: %v", err)
	// 		os.Exit(1)
	// 	}
	// }()
	// go func() {
	// 	for {
	// 		logger.Warn().Msg("[red]Hello, World![white]")
	// 		time.Sleep(250 * time.Millisecond)
	// 	}
	// }()

	miner, err := NewMinerEngine(appContext, manager, &appQuitWG)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create miner engine")
	}

	logger.Info().Str("validator", "******").Str("strategy", cfg.Strategies.Strategy).Msg("⛏️ GOBIUS MINER STARTED! ⛏️")
	//	logger.Info().Str("validator", miner.validator.ValidatorAddress().String()).Str("strategy", cfg.Strategies.Strategy).Msg("⛏️ GOBIUS MINER STARTED! ⛏️")

	err = manager.Start(appQuit)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create start batch manager")
	}

	gpusURLS, ok := cfg.ML.Cog[modelToMine.GetID()]

	if !ok {
		logger.Fatal().Str("model", modelToMine.GetID()).Msg("missing GPU URLs for model - have you set the correct model id?")
	}

	gpus := []*task.GPU{}
	gpuList := zerolog.Arr()
	for id, gpuUrl := range gpusURLS.URL {
		gpuList.Str(gpuUrl)
		gpus = append(gpus, task.NewGPU(id, gpuUrl))
	}

	if !*skipValidation && len(gpus) > 0 {

		validateGpus(modelToMine, gpus, &logger)

		if !miner.services.Engine.VersionCheck() {
			logger.Fatal().Int("minerversion", int(minerVersion.Int64())).Msg("Miner is out of date, please update to the latest version!")
		}
	} else {
		logger.Warn().Msg("skipped model and miner version validation checks!")
	}

	id := len(gpus)
	for i := 0; i < *mockGPUs; i++ {
		gpu := task.NewGPU(id, "")
		gpu.Mock = true
		gpus = append(gpus, gpu)
		id++
	}
	if *mockGPUs > 0 {
		logger.Warn().Msgf("Added %d mock GPUs", *mockGPUs)
	}

	// TODO: move these to a type
	engineContract, err := engine.NewEngine(cfg.BaseConfig.EngineAddress, rpcClient.Client)

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not create engine contract")
	}

	// TODO: move these to a type
	// baseContract, err := basetoken.NewBasetoken(cfg.BaseConfig.BaseTokenAddress, rpcClient.Client)

	// if err != nil {
	// 	logger.Fatal().Err(err).Msg("Could not create engine contract")
	// }

	// TODO: this code only works for websocket/ipc node connections. Add polling support if this fails
	headers := make(chan *types.Header)
	var newHeadSub ethereum.Subscription

	connectToHeaders := func() {
		newHeadSub, err = rpcClient.Client.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to subscribe to new headers")
		}
	}

	connectToHeaders()

	sinkTaskSubmitted := make(chan *engine.EngineTaskSubmitted, 1028)
	var taskEventSub event.Subscription

	// TODO: move these out of here!
	eventAbi, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		panic("error getting engine abi")
	}

	taskSubmittedEvent := eventAbi.Events["TaskSubmitted"].ID

	connectToEvents := func() {

		blockNo, err := rpcClient.Client.BlockNumber(appContext)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to get latest block")
		}

		taskEventSub, err = engineContract.WatchTaskSubmitted(&bind.WatchOpts{
			Start:   &blockNo,
			Context: appContext,
		}, sinkTaskSubmitted, nil, nil, nil)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to subscribe to TaskSubmitted events")
		}
	}

	if cfg.ListenToTaskSubmitted {
		connectToEvents()
	} else {
		logger.Warn().Msg("listening to TaskSubmitted events is disabled")
	}

	// sinkSignalCommitment := make(chan *enginev2.Enginev2SignalCommitment, 1028)
	// signalCommitmentSub, err := engineContract.WatchSignalCommitment(&bind.WatchOpts{
	// 	Start:   &blockNo,
	// 	Context: appContext,
	// }, sinkSignalCommitment, nil, nil)
	// if err != nil {
	// 	logger.Fatal().Err(err).Msg("Failed to subscribe to SignalCommitment events")
	// }
	sinkSolutionSubmitted := make(chan *engine.EngineSolutionSubmitted, 1028)
	var solutionSubmittedSub event.Subscription

	connectToSolutionSubmittedEvents := func() {

		blockNo, err := rpcClient.Client.BlockNumber(appContext)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to get latest block")
		}

		solutionSubmittedSub, err = engineContract.WatchSolutionSubmitted(&bind.WatchOpts{
			Start:   &blockNo,
			Context: appContext,
		}, sinkSolutionSubmitted, nil, nil)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to subscribe to SolutionSubmitted events")
		}
	}

	// httpClient := &http.Client{
	// 	Timeout: time.Second * 10,
	// }
	// ct := geckoterminal.NewClient(httpClient)

	// tokenAddresses := []string{strings.ToLower(cfg.BaseConfig.BaseTokenAddress.String()), WETH_ADDRESS}

	// TODO: uncomment this when we have a way to get token prices
	// prices, err := ct.SimpleTokenPrice("arbitrum-one", tokenAddresses)
	// if err != nil {
	// 	logger.Fatal().Err(err).Msg("error getting token prices")

	// }

	//	logger.Info().Float64("basetokenprice", prices[strings.ToLower(cfg.BaseConfig.BaseTokenAddress.String())]).Float64("eth", prices[WETH_ADDRESS]).Msg("current token prices")

	taskIDs := list.New()
	var mutex sync.Mutex
	newTask := make(chan struct{}, 1)
	var inflightTasks map[string]time.Time = make(map[string]time.Time)
	var erroredTasks, solvedTasks int64 = 0, 0
	var taskHashCache *lru.Cache
	taskHashCache, err = lru.New(1_000_000)
	if err != nil {
		panic(err)
	}

	addTaskToQueue := func(ts *TaskSubmitted) {
		// Push the task ID to the front of the deque
		mutex.Lock()
		taskIDs.PushFront(ts)

		// If the deque is full, remove the oldest task from the back
		if taskIDs.Len() > maxTasks {
			//lastTaskID := taskIDs.Back().Value.([32]byte)
			//log.Println("Task queue is full, ejecting task: ", common.Bytes2Hex(lastTaskID[:]))
			taskIDs.Remove(taskIDs.Back())
		}
		mutex.Unlock()

		//Signal that a new task has been added
		select {
		case newTask <- struct{}{}:
		default:
		}
	}

	//tasksAtStart, _ := miner.services.TaskStorage.TotalTasks()

	// TODO: move this all to the miner engine
	//taskHandler := func(jobid int, event *enginev2.Enginev2TaskSubmitted, gpu task.GPU) {
	taskHandler := func(jobid int, ts *TaskSubmitted, gpu *task.GPU, validateOnly bool) {
		var tx *types.Transaction

		if ts == nil {
			logger.Error().Err(err).Msg("task not set in taskhandler")
			return
		}

		start := time.Now()
		tx, _, err = rpcClient.Client.TransactionByHash(appContext, ts.TxHash)
		elapsed := time.Since(start)
		if err != nil {
			logger.Error().Err(err).Msg("could not get transaction from hash")
			return
		}
		logger.Debug().Str("elapsed", elapsed.String()).Str("hash", ts.TxHash.String()).Msg("got transaction from hash")

		taskIdAsString := ts.TaskId.String()
		start = time.Now()
		_, err = miner.SolveTask(ts.TaskId, tx, gpu, validateOnly)
		elapsed = time.Since(start)

		var solved, errored int64
		if err != nil {
			logger.Error().Err(err).Str("task", taskIdAsString).Msg("solve task failed, adding task back to queue")
			errored = atomic.AddInt64(&erroredTasks, 1)
		} else {
			solved = atomic.AddInt64(&solvedTasks, 1)
		}
		total := solved + errored
		elapsedSeconds := float64(elapsed) / float64(time.Second)
		elapsedString := fmt.Sprintf("%.2f", elapsedSeconds)

		logger.Debug().Int("gpuid", gpu.ID).Str("elapsed", elapsedString).Str("task", taskIdAsString).Int64("total", total).Int64("errors", errored).Msg("task processed")

		// TODO: check who if we didnt mine who did and log freq of the miner
	}

	var jobid int64 = 0

	numWorkers := len(gpus) * cfg.NumWorkersPerGPU

	logger.Info().Int("workerspergpu", cfg.NumWorkersPerGPU).Array("gpus", gpuList).Msgf("running %d workers using the following GPUs", numWorkers)

	appQuitWG.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(workerId int, gpu *task.GPU, wg *sync.WaitGroup) {
			logger.Info().Int("worker", workerId).Int("GPU", gpu.ID).Msg("started worker")
			defer wg.Done()

			// Ticker to periodically check the status of the GPU and re-enable it
			ticker := time.NewTicker(time.Minute * 15) // check every 5 minutes
			defer ticker.Stop()

			for {
				select {
				case <-appQuit.Done():
					logger.Info().Int("worker", workerId).Int("GPU", gpu.ID).Msg("exiting worker")
					return
				case <-ticker.C:
					logger.Info().Int("worker", workerId).Int("GPU", gpu.ID).Msg("resetting gpu state")
					gpu.ResetErrorState()
				case <-newTask:
					if !gpu.IsEnabled() {
						// If the GPU is not enabled, skip this iteration
						continue
					}

					mutex.Lock()
					if taskIDs.Len() == 0 {
						mutex.Unlock()
						continue
					}
					// Pop the latest task from the front of the deque
					element := taskIDs.Front()
					event := element.Value.(*TaskSubmitted)
					taskIDs.Remove(element)
					inflightTasks[event.TaskId.String()] = time.Now()

					sizeOfQ := taskIDs.Len()
					mutex.Unlock()

					// if cfg.Automine.Enabled {

					// 	gpuAvgMiningTime := miner.AvergeGPUTime()
					// 	gpuAvgMiningTime = gpuAvgMiningTime - (850 * time.Millisecond)
					// 	logger.Warn().Int("worker", i).Str("when", gpuAvgMiningTime.String()).Msg("queuing task for automine")

					// 	time.AfterFunc(gpuAvgMiningTime, func() {
					// 		logger.Warn().Int("worker", i).Msg("submitting automine task")
					// 		miner.ProcessAutomine(addTaskToQueue)
					// 	})

					// }

					//.Uint64("eventblock", event.Raw.BlockNumber)
					logger.Debug().Int("worker", workerId).Int("GPU", gpu.ID).Int64("jobid", jobid).Int("pending", sizeOfQ).Msg("starting job")
					taskHandler(workerId, event, gpu, false)
					logger.Debug().Int("worker", workerId).Int("GPU", gpu.ID).Int64("jobid", jobid).Msg("finished job")
					atomic.AddInt64(&jobid, 1)
				}
			}
		}(i, gpus[i/cfg.NumWorkersPerGPU], &appQuitWG)
	}

	getUnsolvedTasks := func(fromBlock int64, blockSize int64) {

		// Get the services from the context
		services, ok := appContext.Value(servicesKey{}).(*Services)
		if !ok {
			log.Fatal("Could not get services from context")
		}

		// Get the current block number
		currentBlock, err := rpcClient.Client.BlockNumber(context.Background())
		if err != nil {
			log.Fatalf("Failed to get current block number: %v", err)
		}

		// Calculate the block number that corresponds to 12 hours ago
		blocksIn12Hours := int64(23800) // 12 hours * 60 minutes/hour * 60 seconds/minute * 4 blocks/second
		toBlock := int64(currentBlock) - blocksIn12Hours
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(0),
			ToBlock:   big.NewInt(0),
			Addresses: []common.Address{cfg.BaseConfig.EngineAddress},
			Topics:    [][]common.Hash{{taskSubmittedEvent}},
		}

		for fromBlock < toBlock {
			nextBlock := fromBlock + blockSize
			if nextBlock > toBlock {
				nextBlock = toBlock
			}

			log.Println("Searching logs ", fromBlock, " thru ", nextBlock)
			query.FromBlock.SetInt64(fromBlock)
			query.ToBlock.SetInt64(nextBlock)

			logs, err := rpcClient.FilterLogs(context.Background(), query)
			if err != nil {
				log.Fatalf("Failed to filter logs: %v", err)
			}

			if len(logs) > 0 {
				log.Println("found tasksubmitted events:", len(logs))
			}

			for _, currentlog := range logs {
				parsedLog, err := engineContract.ParseTaskSubmitted(currentlog)
				if err != nil {
					log.Fatalf("Failed to filter logs: %v", err)
				}

				sol, err := engineContract.Solutions(nil, parsedLog.Id)
				if err != nil {
					log.Fatalf("Failed to filter logs: %v", err)
				}
				if sol.Blocktime == 0 {
					taskId := task.TaskId(parsedLog.Id).String()
					log.Println("Add unsolved task to queueueue:", taskId)
					services.TaskStorage.AddTask(taskId, currentlog.TxHash.String())
				}
			}

			fromBlock = nextBlock + 1
		}
	}

	if *taskScanner > 0 {
		go getUnsolvedTasks(int64(*taskScanner), 1000)
	}

	currentBlockNumer := uint64(0)

	switch cfg.Strategies.Strategy {

	case "bulkmine":
		go func() {

			addTask := func(task *TaskSubmitted) bool {
				mutex.Lock()
				taskIDs.PushFront(task)
				mutex.Unlock()

				select {
				case newTask <- struct{}{}:
				case <-appQuit.Done():
					logger.Info().Msg("shutting down task queue worker")
					return true
				}

				return false
			}
			poppedCount := 0

			for {

				//logger.Info().Msg("waiting for tasks to be added to storage queue!")

				taskId, txHash, err := miner.services.TaskStorage.PopTask()
				if err != nil {
					if err == sql.ErrNoRows {
						logger.Warn().Msgf("queue is empty, nothing to do, will sleep. popped so far: %d", poppedCount)
						time.Sleep(2 * time.Second)
						// The list is empty, continue to the next iteration
						continue
					}
					logger.Error().Err(err).Msgf("could not pop task from storage, will sleep. popped so far: %d", poppedCount)
					time.Sleep(5 * time.Second)
					continue
				}

				_, found := taskHashCache.Get(taskId)
				if found {
					logger.Error().Msgf("trying to pop same task again: %s", taskId.String())
					continue
				}

				taskHashCache.Add(taskId, struct{}{})

				ts := &TaskSubmitted{
					TaskId: taskId,
					TxHash: txHash,
				}
				poppedCount++

				logger.Info().Msgf("popped task: %s", taskId.String())

				if addTask(ts) {
					return
				}
				time.Sleep(10 * time.Millisecond)
			}
		}()

	case "solutionsampler":

		logger.Info().Msg("strategy is sample other solutions!")

		connectToSolutionSubmittedEvents()

		solutionSampler := func() {
			ticker := time.NewTicker(time.Minute)
			defer ticker.Stop()

			maxTaskSampleSize := numWorkers
			tasksSamples := make([]*TaskSubmitted, 0, maxTaskSampleSize)
			sampleIndex := 0
			for {
				select {
				case <-appQuit.Done():
					logger.Info().Msg("shutting down sampler")
					return
				case err := <-solutionSubmittedSub.Err():
					// TODO: don't make this a fatal error - have some soft recovery
					log.Printf("Error from solutionSubmittedSub: %v", err)
					time.Sleep(time.Second * 5)
					connectToSolutionSubmittedEvents()

				case event := <-sinkSolutionSubmitted:

					taskHash, found := taskHashCache.Get(event.Task)
					if !found {
						continue
					}
					ts := &TaskSubmitted{
						TaskId: event.Task,
						TxHash: taskHash.(common.Hash),
					}

					sampleIndex++
					if len(tasksSamples) < maxTaskSampleSize {
						tasksSamples = append(tasksSamples, ts)
					} else {
						j := rand.Intn(sampleIndex)
						if j < maxTaskSampleSize {
							tasksSamples[j] = ts
						}
					}
				case <-ticker.C:
					var wgWorkers sync.WaitGroup

					logger.Info().Int("samples", len(tasksSamples)).Msg("enough solution samples taken, processing...")

					for i := 0; i < len(tasksSamples); i++ {
						wgWorkers.Add(1)
						go func(ts *TaskSubmitted, workerId int, gpu *task.GPU, wg *sync.WaitGroup) {
							defer wgWorkers.Done()
							taskIdAsString := ts.TaskId.String()

							logger.Info().Int("workerid", workerId).Str("task", taskIdAsString).Msg("processing sampled task")

							tx, _, err := rpcClient.Client.TransactionByHash(appContext, ts.TxHash)
							if err != nil {
								logger.Error().Err(err).Msg("could not get transaction from hash")
								return
							}
							cid, err := miner.SolveTask(ts.TaskId, tx, gpu, true)
							if err != nil {
								logger.Error().Err(err).Msg("solve failed")
								return
							}
							if cid == nil {
								logger.Info().Msg("solve didn't return a cid for task")
								return
							}

							res, err := engineContract.Solutions(nil, ts.TaskId)
							if err != nil {
								logger.Err(err).Msg("error getting solution information")
								return
							}

							if res.Blocktime > 0 {

								solversCid := common.Bytes2Hex(res.Cid[:])
								ourCid := common.Bytes2Hex(cid)
								logger.Info().Msgf("checking cids: %s and %s", ourCid, solversCid)

								if ourCid != solversCid {
									logger.Warn().Msg("=======================================================================")
									logger.Warn().Msg("  WARNING: our solution cid does not match the solvers cid!")
									logger.Warn().Msg("  our cid: " + ourCid)
									logger.Warn().Msg("  ther cid: " + solversCid)
									logger.Warn().Str("validator", res.Validator.String()).Msg("  solvers address")
									logger.Warn().Msg("========================================================================")
								}
							}

						}(tasksSamples[i], i, gpus[i/cfg.NumWorkersPerGPU], &wgWorkers)
					}
					wgWorkers.Wait()
					logger.Info().Msg("GPUS FINISHED NOW")
					tasksSamples = tasksSamples[:0] // reset samples
					sampleIndex = 0
				}
			}
		}

		go solutionSampler()
	}

	if cfg.ListenToTaskSubmitted {

		go func() {
			maxBackoff := time.Second * 30
			currentBackoff := time.Second
			for {
				select {
				case err := <-taskEventSub.Err():
					log.Printf("Error from taskEventSub: %v - redialling in: %s\n", err, currentBackoff.String())

					time.Sleep(currentBackoff)
					currentBackoff *= 2
					currentBackoff += time.Duration(rand.Intn(1000)) * time.Millisecond
					if currentBackoff > maxBackoff {
						currentBackoff = maxBackoff
					}

					connectToEvents()
				// case <-appContext.Done():
				// 	logger.Info().Msg("Shutting down")
				// 	return
				case event := <-sinkTaskSubmitted:
					taskId := task.TaskId(event.Id)

					ts := &TaskSubmitted{
						TaskId: taskId,
						TxHash: event.Raw.TxHash,
					}

					switch cfg.Strategies.Strategy {
					case "task":
						kick := event.Raw.BlockNumber < currentBlockNumer-1
						logger.Debug().Str("taskid", taskId.String()).Uint64("eventblock", event.Raw.BlockNumber).Uint64("currentblock", currentBlockNumer).Bool("eject?", kick).Msg("[Task Submitted]")

						if !cfg.StealTasks && kick {
							continue
						}
						addTaskToQueue(ts)
					case "solutionsampler":
						taskHashCache.Add(event.Id, event.Raw.TxHash)
					}
				}
			}
		}()
	}

	maxBackoffHeader := time.Second * 30
	currentBackoffHeader := time.Second

	for {

		select {

		case <-appQuit.Done():
			logger.Info().Msg("shutting down main loop")
			goto exit_app
		case err := <-newHeadSub.Err():
			if err == nil {
				continue
			}
			log.Printf("Error from newHeadSub: %v - redialling in: %s\n", err, currentBackoffHeader.String())
			newHeadSub.Unsubscribe()

			time.Sleep(currentBackoffHeader)
			currentBackoffHeader *= 2
			currentBackoffHeader += time.Duration(rand.Intn(1000)) * time.Millisecond
			if currentBackoffHeader > maxBackoffHeader {
				currentBackoffHeader = maxBackoffHeader
			}

			connectToHeaders()

		case h := <-headers:
			//blockTime := time.Unix(int64(h.Time), 0)

			// /log.Println("New block: ", h.Number.String(), blockTime) // print the block number
			//basefeeStr := "not avail"
			if h.BaseFee != nil {
				// /basefeeStr = h.BaseFee.String()

				// update basefee
				rpcClient.SetBaseFee(h.BaseFee)
			}
			// TODO: atomic updates pls
			currentBlockNumer = h.Number.Uint64()
		}
	}
exit_app:
	logger.Info().Msg("Waiting for application workers to finish")
	appQuitWG.Wait()
	logger.Info().Msg("bye! 👋")
}
