package main

import (
	"container/list"
	"context"
	"database/sql"
	"embed"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"gobius/bindings/engine"
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
	"sync"
	"sync/atomic"
	"time"

	"gobius/metrics"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pressly/goose/v3"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

//go:embed sql/sqlite/migrations/*.sql
var embedMigrations embed.FS

type TaskSubmitted struct {
	TaskId task.TaskId
	TxHash common.Hash
}

const (
	maxTasks                       = 50 // Maximum number of tasks to store
	maxExitAttempts                = 3  // Maximum number of attempts to exit the application
	appVersion                     = "0.0.1"
	taskSubmittedChannelBufferSize = 1024
	appName                        = `
   ┏┓┏┓┳┓┳┳┳┏┓     
   ┃┓┃┃┣┫┃┃┃┗┓    
   ┗┛┗┛┻┛┻┗┛┗┛    
`
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

	logger = zerolog.New(multi).Level(level).With().Timestamp().Logger()

	return fileLogger, logger
}

func setupCloseHandler(logger *zerolog.Logger) (ctx context.Context, cancel context.CancelFunc) {
	// Create a cancellation context.
	ctx, cancel = context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		exitCounter := maxExitAttempts

		for {
			<-c

			switch {
			case exitCounter < 3 && exitCounter > 0:
				logger.Error().Int("attempt", maxExitAttempts+1-exitCounter).Msgf("already shutting down. will force exit after %d more attempt(s).", exitCounter)
				exitCounter--
			case exitCounter == 0:
				logger.Error().Msgf("ctrl+C pressed %d times. forcing exit.", maxExitAttempts)
				os.Exit(0)
			default:
				logger.Info().Msg("ctrl+C detected. stopping process...")
				cancel()
				exitCounter--
			}
		}
	}()

	return ctx, cancel
}

// type TaskManagerI interface {
// 	GetCurrentReward() (*big.Int, error)
// 	GetTotalTasks() int64
// 	GetClaims() int64
// 	GetSolutions() int64
// 	GetCommitments() int64
// }

// TODO: put this into base config ?
var minerVersion = big.NewInt(5)

type Miner struct {
	services  *Services
	ctx       context.Context
	mu        sync.Mutex
	gpura     *utils.RunningAverage
	validator IValidator
	wg        *sync.WaitGroup
	// gpus will store GPU instances for metrics/TUI
	gpus []*task.GPU
	// TUI related fields
	p *tea.Program
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

// connects miner account to the engine contract
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

func (m *Miner) GetGPUInfo() []metrics.GPUInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	var gpuInfos []metrics.GPUInfo
	for _, gpu := range m.gpus {
		status := "Mining"
		if !gpu.IsEnabled() {
			status = "Error"
		}

		gpuInfos = append(gpuInfos, metrics.GPUInfo{
			ID:        gpu.ID,
			Status:    status,
			SolveTime: m.gpura.Average(),
		})
	}
	return gpuInfos
}

func (m *Miner) SolveTask(taskId task.TaskId, tx *types.Transaction, gpu *task.GPU, validateOnly bool) ([]byte, error) {

	taskIdStr := taskId.String()

	taskData := m.services.AutoMineParams

	inputRaw := string(taskData.Input)

	var result map[string]interface{}
	err := json.Unmarshal(taskData.Input, &result)

	if err != nil {
		m.services.Logger.Error().Err(err).Msg("could not unmarshal input")
		return nil, err
	}

	m.services.Logger.Debug().Str("input", inputRaw).Msg("decoded information")

	taskInfo, err := m.services.Engine.LookupTask(taskId)
	if err != nil {
		m.services.Logger.Error().Err(err).Msg("could not lookup task")
		return nil, err
	}

	modelId := common.Bytes2Hex(taskInfo.Model[:])
	model := models.ModelRegistry.GetModel(modelId)
	if model == nil {
		m.services.Logger.Error().Str("model", modelId).Err(err).Msg("could not find model")
		return nil, err
	}

	hydrated, err := model.HydrateInput(result, taskId.TaskId2Seed())

	if err != nil {
		m.services.Logger.Error().Err(err).Msg("could not hydrate input")
		return nil, err
	}

	output, err := json.Marshal(hydrated)
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
			data, err := gpu.GetMockCid(taskIdStr, hydrated)
			if err != nil {
				return nil, err
			}
			cid, err = ipfs.GetIPFSHashFast(data)
			if err != nil {
				return nil, err
			}
		} else {
			cid, err = model.GetCID(gpu, taskIdStr, hydrated)
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

	// Use a separate goroutine without WaitGroup tracking for solution submission
	go m.validator.SubmitSolution(validator, taskId, cid)

	return cid, nil
}

func average(times []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, time := range times {
		total += time
	}
	return total / time.Duration(len(times))
}

func validateGpus(model models.ModelInterface, gpus []*task.GPU, logger *zerolog.Logger) error {

	//findFastestPrompt(modelToMine, gpus, &logger)

	logger.Info().Str("model", model.GetID()).Msg("validating model on gpu(s)")

	mu := sync.Mutex{}
	fastestTimeSeen := time.Duration(0)

	var wg sync.WaitGroup
	times := make([]time.Duration, len(gpus))

	for i, gpu := range gpus {
		wg.Add(1)

		go func(i int, gpu *task.GPU) {
			defer wg.Done()

			start := time.Now()

			err := model.Validate(gpu, "startup-test-taskid")

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

// basic log output router to direct logging to either console or tui if avail
type logrouter struct {
	view     *tui.CustomTextView
	Headless atomic.Bool
	writer   io.Writer
}

func (tw *logrouter) SetView(view *tui.CustomTextView) {
	tw.view = view
	tw.writer = tview.ANSIWriter(view)
}

func (tw *logrouter) Write(p []byte) (n int, err error) {
	isHeadless := tw.Headless.Load()
	if tw.view == nil || isHeadless {
		return os.Stderr.Write(p)
	} else {
		//return tw.writer.Write(p)
		//return .Write(p)q
		return fmt.Fprintf(tw.writer, "%s", p)
	}
}

func main() {
	var appQuitWG sync.WaitGroup

	configPath := flag.String("config", "config.json", "Path to the configuration file")
	skipValidation := flag.Bool("skipvalidation", false, "Skip safety checks and validation of the model and miner version")
	logLevel := flag.Int("loglevel", 1, "Set the logging level")
	testnetType := flag.Int("testnet", 0, "Run using testnet - 1 = local, 2 = nova testnet")
	mockGPUs := flag.Int("mockgpus", 0, "mock gpus for testing")
	taskScanner := flag.Int("taskscanner", 0, "scan blocks for unsolved tasks")
	headless := flag.Bool("headless", true, "Run in headless mode without the dashboard UI")

	flag.Parse()

	// Create a set of the flags that were set
	setFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		setFlags[f.Name] = true
	})

	fmt.Println(appName)
	fmt.Printf("Version: %s\n\n", appVersion)

	cfg, err := config.InitAppConfig(*configPath, *testnetType)
	if err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	// Check if logLevel was set
	if !setFlags["loglevel"] {
		// If logLevel was not set on the command line, use the value from the config
		*logLevel = int(cfg.LogLevel)
	}

	// logWriter := &tui.LogViewWriter{
	// 	App:      nil,
	// 	Headless: *headless,
	// }

	logWriter := &logrouter{
		view:     nil,
		Headless: atomic.Bool{},
	}
	logWriter.Headless.Store(*headless)

	logFile, logger := initLogging(cfg.LogPath, zerolog.Level(*logLevel), logWriter)

	defer logFile.Close()

	if cfg.EvilMode {
		logger.Warn().Msg("TESTNET MODE ENABLED - DO NOT USE ON MAINNET")
	}

	appQuit, appCancel := setupCloseHandler(&logger)

	rpcClient, err := client.NewClient(cfg.Blockchain.RPCURL, appQuit, cfg.Blockchain.EthersGas, cfg.Blockchain.BasefeeX, cfg.Blockchain.ForceGas, cfg.Blockchain.GasOverride)

	if err != nil {
		logger.Error().Err(err).Msgf("error connecting to RPC: %s", cfg.Blockchain.RPCURL)
		return
	}

	txRpcClient := rpcClient
	if cfg.Blockchain.SenderRPCURL != "" {
		txRpcClient, err = client.NewClient(cfg.Blockchain.SenderRPCURL, appQuit, cfg.Blockchain.EthersGas, cfg.Blockchain.BasefeeX, cfg.Blockchain.ForceGas, cfg.Blockchain.GasOverride)

		if err != nil {
			logger.Error().Err(err).Msgf("error connecting to sender RPC: %s", cfg.Blockchain.SenderRPCURL)
			return
		}
	}

	var clients []*client.Client
	for _, curl := range cfg.Blockchain.ClientRPCURLs {
		c, err := client.NewClient(curl, appQuit, cfg.Blockchain.EthersGas, cfg.Blockchain.BasefeeX, cfg.Blockchain.ForceGas, cfg.Blockchain.GasOverride)

		if err != nil {
			logger.Error().Err(err).Msgf("error connecting to client RPC: %s", curl)
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

	var ipfsClient ipfs.IPFSClient

	switch cfg.IPFS.Strategy {
	case "mock":
		ipfsClient, err = ipfs.NewMockIPFSClient(*cfg, true)
	case "http_client":
		ipfsClient, err = ipfs.NewHttpIPFSClient(*cfg, true)
	default:
		logger.Fatal().Str("strategy", cfg.IPFS.Strategy).Msg("invalid IPFS strategy")
	}

	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to IPFS")
	}

	models.InitModelRegistry(ipfsClient, cfg, &logger)

	modelToMine := models.ModelRegistry.GetModel(cfg.Strategies.Model)
	if modelToMine == nil {
		logger.Error().Str("model", cfg.Strategies.Model).Msg("model specified in config was not found in enabled models")
		return
	}

	miner, err := NewMinerEngine(appContext, manager, &appQuitWG)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create miner engine")
	}

	err = manager.Start(appQuit)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create start batch manager")
	}

	logger.Info().Str("strategy", cfg.Strategies.Strategy).Msg("⛏️ GOBIUS MINER STARTED! ⛏️")
	//	logger.Info().Str("validator", miner.validator.ValidatorAddress().String()).Str("strategy", cfg.Strategies.Strategy).Msg("⛏️ GOBIUS MINER STARTED! ⛏️")

	gpusURLS, ok := cfg.ML.Cog[modelToMine.GetID()]

	if !ok {
		logger.Fatal().Str("model_id", modelToMine.GetID()).Msg("missing GPU URLs for model - have you set the correct model id?")
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
			logger.Fatal().Int("minerversion", int(minerVersion.Int64())).Msg("miner is out of date, please update to the latest version!")
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
		logger.Warn().Msgf("added %d mock GPUs", *mockGPUs)
	}

	miner.gpus = gpus

	// TODO: move these to a type
	engineContract, err := engine.NewEngine(cfg.BaseConfig.EngineAddress, rpcClient.Client)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not create engine contract")
	}

	// TODO: this code only works for websocket/ipc node connections. Add polling support if this fails
	headers := make(chan *types.Header)
	var newHeadSub ethereum.Subscription

	connectToHeaders := func() {
		ctx, cancel := context.WithTimeout(appContext, 5*time.Second)
		defer cancel()

		newHeadSub, err = rpcClient.Client.SubscribeNewHead(ctx, headers)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to subscribe to new headers, RPC must be websocket/ipc only, not http(s)")
		}
	}

	connectToHeaders()

	sinkTaskSubmitted := make(chan *engine.EngineTaskSubmitted, taskSubmittedChannelBufferSize)
	var taskEventSub event.Subscription

	// TODO: move these out of here!
	eventAbi, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		panic("error getting engine abi")
	}

	taskSubmittedEvent := eventAbi.Events["TaskSubmitted"].ID

	connectToEvents := func() {

		ctx, cancel := context.WithTimeout(appContext, 5*time.Second)
		defer cancel()

		blockNo, err := rpcClient.Client.BlockNumber(appContext)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to get latest block")
		}

		taskEventSub, err = engineContract.WatchTaskSubmitted(&bind.WatchOpts{
			Start:   &blockNo,
			Context: ctx,
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

	dashboard := tui.NewDashboard()

	if !*headless {
		// Update our log router to use the logviewer
		logWriter.SetView(dashboard.LogViewer.CustomTextView)
		go func() {

			dashboard.Run()
			// disable writing to log viwer on exit
			// we don't set the view to nil as this is unsafe instead we use a atomic bool
			logWriter.Headless.Store(true)
			appCancel()
			// ticker := time.NewTicker(time.Second)
			// defer ticker.Stop()

			// for {
			// 	select {
			// 	case <-dashboard.GetQuitSignal():
			// 		appCancel() // Trigger app shutdown using the cancel function
			// 		return
			// 	case <-ticker.C:

			// 	}
			// }
		}()

	} else {
		logger.Info().Msg("running in headless mode; dashboard disabled")
	}
	/*var p *tea.Program
	dashboard := tui.NewDashboard()
	if !*headless {
		p = tea.NewProgram(
			dashboard,
			tea.WithAltScreen(),
			tea.WithMouseAllMotion(),
		)

		// Set the program in the log writer
		logWriter.App = p

		// Set the writer in the dashboard
		dashboard.SetLogWriter(logWriter)

		// Start dashboard updater
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-dashboard.GetQuitSignal():
					appCancel() // Trigger app shutdown using the cancel function
					return
				case <-ticker.C:
					// Update validator metrics
					validatorMetrics := tui.ValidatorMetrics{
						SessionTime:      manager.GetSessionTime(),
						SolvedLastMinute: manager.GetSolvedLastMinute(),
						SolutionsLastMinute: struct {
							Success int64
							Total   int64
							Rate    float64
						}{
							Success: manager.GetSuccessCount(),
							Total:   manager.GetTotalCount(),
							Rate:    manager.GetSuccessRate(),
						},
						AverageSolutionRate:    manager.GetAverageSolutionRate(),
						AverageSolutionsPerMin: manager.GetAverageSolutionsPerMin(),
						AverageSolvesPerMin:    manager.GetAverageSolvesPerMin(),
					}
					p.Send(dashboard.UpdateValidatorMetrics(validatorMetrics))

					// Update financial metrics
					financialMetrics := tui.FinancialMetrics{
						TokenIncomePerMin:  manager.GetTokenIncomePerMin(),
						TokenIncomePerHour: manager.GetTokenIncomePerHour(),
						TokenIncomePerDay:  manager.GetTokenIncomePerDay(),
						IncomePerMin:       manager.GetIncomePerMin(),
						IncomePerHour:      manager.GetIncomePerHour(),
						IncomePerDay:       manager.GetIncomePerDay(),
						ProfitPerMin:       manager.GetProfitPerMin(),
						ProfitPerHour:      manager.GetProfitPerHour(),
						ProfitPerDay:       manager.GetProfitPerDay(),
					}
					p.Send(dashboard.UpdateFinancialMetrics(financialMetrics))

					// Update GPU metrics
					dashboard.SendGPUMetrics(p, miner.GetGPUInfo())
				}
			}
		}()

		go func() {
			if _, err := p.Run(); err != nil {
				logger.Fatal().Err(err).Msg("failed to start dashboard")
			}
		}()

	} else {
		logger.Info().Msg("running in headless mode - dashboard disabled")
	}*/

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

			// Move defer after the logger so we don't lose shutdown messages
			defer func() {
				logger.Info().Int("worker", workerId).Int("GPU", gpu.ID).Msg("worker finished")
				wg.Done()
			}()

			// Ticker to periodically check the status of the GPU and re-enable it
			ticker := time.NewTicker(time.Minute * 15) // check every 5 minutes
			defer ticker.Stop()

			for {
				select {
				case <-appQuit.Done():
					logger.Info().Int("worker", workerId).Int("GPU", gpu.ID).Msg("exiting worker")
					return
				case <-ticker.C:
					// Check if we're shutting down before doing work
					if appQuit.Err() != nil {
						return
					}
					logger.Info().Int("worker", workerId).Int("GPU", gpu.ID).Msg("resetting gpu state")
					gpu.ResetErrorState()
				case <-newTask:
					// Check if we're shutting down before doing work
					if appQuit.Err() != nil {
						return
					}
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
					log.Println("Add unsolved task to queueue:", taskId)
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
		appQuitWG.Add(1)
		go func() {
			defer appQuitWG.Done()
			addTask := func(task *TaskSubmitted) bool {
				// Don't add tasks if we're shutting down
				if appQuit.Err() != nil {
					return true
				}

				mutex.Lock()
				taskIDs.PushFront(task)
				mutex.Unlock()

				select {
				case newTask <- struct{}{}:
					return false
				case <-appQuit.Done():
					logger.Info().Msg("shutting down task queue worker")
					return true
				}
			}
			poppedCount := 0

			for {
				select {
				case <-appQuit.Done():
					logger.Info().Msg("bulkmine worker shutting down")
					return
				default:
					if appQuit.Err() != nil {
						return
					}

					taskId, txHash, err := miner.services.TaskStorage.PopTask()
					if err != nil {
						if err == sql.ErrNoRows {
							logger.Warn().Msgf("queue is empty, nothing to do, will sleep. popped so far: %d", poppedCount)
							time.Sleep(2 * time.Second)
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
			}
		}()

	case "solutionsampler":
		logger.Info().Msg("strategy is sample other solutions!")

		connectToSolutionSubmittedEvents()

		appQuitWG.Add(1)
		go func() {
			defer func() {
				logger.Info().Msg("solution sampler shutting down")
				appQuitWG.Done()
			}()

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
						logger.Warn().Msgf("solution submitted sub error: %v, will retry in 5s", err)

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
						if appQuit.Err() != nil {
							return
						}

						var wgWorkers sync.WaitGroup
						logger.Info().Int("samples", len(tasksSamples)).Msg("enough solution samples taken, processing...")

						for i := 0; i < len(tasksSamples); i++ {
							if appQuit.Err() != nil {
								wgWorkers.Wait()
								return
							}
							wgWorkers.Add(1)
							go func(ts *TaskSubmitted, workerId int, gpu *task.GPU, wg *sync.WaitGroup) {
								defer wg.Done()
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
			solutionSampler()
		}()
	}

	if cfg.ListenToTaskSubmitted {
		appQuitWG.Add(1)
		go func() {
			defer appQuitWG.Done()
			maxBackoff := time.Second * 30
			currentBackoff := time.Second
			for {
				select {
				case <-appQuit.Done():
					return
				case err := <-taskEventSub.Err():
					logger.Warn().Msgf("task event sub error: %v, will retry in %s", err, currentBackoff.String())

					time.Sleep(currentBackoff)
					currentBackoff *= 2
					currentBackoff += time.Duration(rand.Intn(1000)) * time.Millisecond
					if currentBackoff > maxBackoff {
						currentBackoff = maxBackoff
					}

					connectToEvents()
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
			logger.Warn().Msgf("new head sub error: %v, will retry in %s", err, currentBackoffHeader.String())
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
	logger.Info().Msg("waiting for application workers to finish")
	// Wait for all workers to finish
	appQuitWG.Wait()

	// Now that all tasks are complete, properly cleanup the TUI
	// if !*headless && p != nil {
	// 	time.Sleep(100 * time.Millisecond) // Give a moment for the final frame to render
	// 	p.Quit()
	// 	p.Wait()
	// }

	logger.Info().Msg("bye! 👋")
}
