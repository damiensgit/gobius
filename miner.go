package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/hex"
	"encoding/json"
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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Embed section
// this stores the migrations for the database
//
//go:embed sql/sqlite/migrations/*.sql
var embedMigrations embed.FS

// Constants section
const (
	maxTasks                       = 50 // Maximum number of tasks to store
	maxExitAttempts                = 3  // Maximum number of attempts to exit the application
	appVersionMajor                = 1
	appVersionMinor                = 0
	appVersionPatch                = 0
	taskSubmittedChannelBufferSize = 1024
	appName                        = `
   ┏┓┏┓┳┓┳┳┳┏┓     
   ┃┓┃┃┣┫┃┃┃┗┓    v%d.%02d.%02d / engine v%d
   ┗┛┗┛┻┛┻┗┛┗┛    

`
)

// Variables section
// TODO: put this into base config ?
var minerEngineVersion = big.NewInt(5)

// Types section
type Miner struct {
	services  *Services
	validator IValidator
}

// zerologAdapter adapts a zerolog.Logger to satisfy goose.Logger interface
type zerologAdapter struct {
	logger zerolog.Logger
}

func (z *zerologAdapter) Printf(format string, v ...interface{}) {
	z.logger.Info().Msgf(format, v...)
}

func (z *zerologAdapter) Fatalf(format string, v ...interface{}) {
	z.logger.Fatal().Msgf(format, v...)
}

func initLogging(file string, level zerolog.Level) (logFile io.WriteCloser, logger zerolog.Logger) {
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.TimeFieldFormat = time.RFC3339Nano

	fileLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     5, // days
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000000000"}
	consoleWriter.FormatTimestamp = utils.ZerologConsoleFormatTimestamp(consoleWriter.TimeFormat)
	multi := zerolog.MultiLevelWriter(consoleWriter, fileLogger)
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
				logger.Error().Msgf("ctrl+c pressed %d times. forcing exit.", maxExitAttempts)
				os.Exit(0)
			default:
				logger.Info().Msg("ctrl+c detected. stopping process...")
				cancel()
				exitCounter--
			}
		}
	}()

	return ctx, cancel
}

// connects miner account to the engine contract
func NewMinerEngine(services *Services, validator IValidator, wg *sync.WaitGroup) (*Miner, error) {

	miner := &Miner{
		services:  services,
		validator: validator,
	}

	return miner, nil
}

func (m *Miner) SolveTask(ctx context.Context, taskId task.TaskId, tx *types.Transaction, gpu *task.GPU, validateOnly bool) ([]byte, error) {

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
		//start := time.Now()
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
		//elapsed := time.Since(start)
		//m.gpura.Add(elapsed)
		if err != nil {
			m.services.Logger.Error().Err(err).Msg("error on gpu, incrementing error counter")
			gpu.IncrementErrorCount()
			return nil, err
		}
		//m.services.Logger.Debug().Str("cid", "0x"+hex.EncodeToString(cid)).Str("elapsed", elapsed.String()).Str("average", m.gpura.Average().String()).Msg("gpu finished & returned result")

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

	err = m.validator.SubmitIpfsCid(validator, taskId, cid)
	if err != nil {
		m.services.Logger.Warn().Err(err).Msg("ipfs cid submission failed")
	}

	// Use a separate goroutine without WaitGroup tracking for solution submission
	go m.validator.SubmitSolution(validator, taskId, cid)

	return cid, nil
}

func main() {
	var appQuitWG sync.WaitGroup

	configPath := flag.String("config", "config.json", "Path to the configuration file")
	skipValidation := flag.Bool("skipvalidation", false, "Skip safety checks and validation of the model and miner version")
	logLevel := flag.Int("loglevel", 1, "Set the logging level")
	testnetType := flag.Int("testnet", 0, "Run using specified testnet: 1 = local, 2 = arbitrum sepolia testnet")
	mockGPUs := flag.Int("mockgpus", 0, "mock gpus for testing")
	//taskScanner := flag.Int("taskscanner", 0, "scan blocks for unsolved tasks")
	headless := flag.Bool("headless", true, "Run in headless mode without the dashboard UI")

	flag.Parse()

	// Create a set of the flags that were set
	setFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		setFlags[f.Name] = true
	})

	fmt.Printf(appName, appVersionMajor, appVersionMinor, appVersionPatch, minerEngineVersion)

	cfg, err := config.InitAppConfig(*configPath, *testnetType)
	if err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	// Check if logLevel was set
	if !setFlags["loglevel"] {
		// If logLevel was not set on the command line, use the value from the config
		*logLevel = int(cfg.LogLevel)
	}

	logFile, logger := initLogging(cfg.LogPath, zerolog.Level(*logLevel))
	defer logFile.Close()

	if cfg.EvilMode {
		logger.Warn().Msg("TESTNET MODE ENABLED - DO NOT USE ON MAINNET")
	}

	appQuit, appCancel := setupCloseHandler(&logger)

	rpcClient, err := client.NewClient(cfg.Blockchain.RPCURL, appQuit, cfg.Blockchain.EthersGas, cfg.Blockchain.BasefeeX, cfg.Blockchain.ForceGas, cfg.Blockchain.GasOverride)

	if err != nil {
		logger.Fatal().Err(err).Msgf("error connecting to RPC: %s", cfg.Blockchain.RPCURL)
	}

	txRpcClient := rpcClient
	if cfg.Blockchain.SenderRPCURL != "" {
		txRpcClient, err = client.NewClient(cfg.Blockchain.SenderRPCURL, appQuit, cfg.Blockchain.EthersGas, cfg.Blockchain.BasefeeX, cfg.Blockchain.ForceGas, cfg.Blockchain.GasOverride)

		if err != nil {
			logger.Fatal().Err(err).Msgf("error connecting to sender RPC: %s", cfg.Blockchain.SenderRPCURL)
		}
	}

	var clients []*client.Client
	for _, curl := range cfg.Blockchain.ClientRPCURLs {
		c, err := client.NewClient(curl, appQuit, cfg.Blockchain.EthersGas, cfg.Blockchain.BasefeeX, cfg.Blockchain.ForceGas, cfg.Blockchain.GasOverride)

		if err != nil {
			logger.Fatal().Err(err).Msgf("error connecting to client RPC: %s", curl)
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
		logger.Fatal().Err(err).Msg("error setting pragma journal mode on sqlite")
	}

	goose.SetBaseFS(embedMigrations)

	goose.SetLogger(&zerologAdapter{logger: logger})

	if err := goose.SetDialect("sqlite3"); err != nil {
		logger.Fatal().Err(err).Msg("error setting goose dialect to sqlite3")
	}

	if err := goose.Up(sqlite, "sql/sqlite/migrations"); err != nil {
		logger.Fatal().Err(err).Msg("database migration error")
	}

	var ipfsOracle ipfs.OracleClient

	if cfg.IPFS.OracleURL != "" {
		timeout, err := time.ParseDuration(cfg.IPFS.Timeout)
		if err != nil {
			logger.Fatal().Err(err).Msg("invalid IPFS oracle timeout")
		}
		ipfsOracle = ipfs.NewHTTPOracleClient(cfg.IPFS.OracleURL, timeout)
	} else {
		ipfsOracle = ipfs.NewMockOracleClient()
	}

	appContext, appServices, err := NewApplicationContext(rpcClient, txRpcClient, clients, sqlite, logger, cfg, ipfsOracle, context.Background(), appQuit)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not create application context")
	}

	manager, err := NewBatchTransactionManager(appServices, appContext, &appQuitWG)

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

	models.InitModelRegistry(ipfsClient, cfg, logger)
	modelToMine := models.ModelRegistry.GetModel(cfg.Strategies.Model)
	if modelToMine == nil {
		logger.Fatal().Str("model", cfg.Strategies.Model).Msg("model specified in config was not found in enabled models")
	}

	miner, err := NewMinerEngine(appServices, manager, &appQuitWG)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create miner engine")
	}

	logger.Info().Str("strategy", cfg.Strategies.Strategy).Msg("⛏️ GOBIUS MINER STARTED! ⛏️")
	//	logger.Info().Str("validator", miner.validator.ValidatorAddress().String()).Str("strategy", cfg.Strategies.Strategy).Msg("⛏️ GOBIUS MINER STARTED! ⛏️")

	// GPU Pool
	gpuPool, err := NewGPUPool(cfg, logger, modelToMine.GetID(), *mockGPUs)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create GPU pool")
	}

	if !*skipValidation {
		err = gpuPool.ValidateGPUs(modelToMine)
		if err != nil {
			// Log the error, but potentially continue if some GPUs are okay?
			// For now, treat validation failure as fatal if not skipped.
			logger.Fatal().Err(err).Msg("GPU validation failed")
		}
		if !miner.services.Engine.VersionCheck(minerEngineVersion) {
			logger.Fatal().Int("minerversion", int(minerEngineVersion.Int64())).Msg("miner is out of date, please update!")
		}
		logger.Info().Msg("GPU validation and engine version check passed")
	} else {
		logger.Warn().Msg("skipped model validation and engine version checks!")
	}

	err = manager.Start(appQuit)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create start batch manager")
	}

	// TODO: this code only works for websocket/ipc node connections. Add polling support if this fails
	headers := make(chan *types.Header)
	var newHeadSub ethereum.Subscription
	connectToHeaders := func() {
		ctx, cancel := context.WithTimeout(appContext, 5*time.Second)
		defer cancel()

		newHeadSub, err = rpcClient.Client.SubscribeNewHead(ctx, headers)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to subscribe to new headers, tip: RPC must be websocket/ipc only, not http(s)")
		}
		logger.Info().Msg("subscribed to new headers")
	}

	connectToHeaders()

	engineContract, err := engine.NewEngine(cfg.BaseConfig.EngineAddress, rpcClient.Client)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create engine contract")
	}

	// TaskSubmitted Event Subscription
	sinkTaskSubmitted := make(chan *engine.EngineTaskSubmitted, 1024) // Buffer size
	var taskEventSub event.Subscription
	connectToTaskEvents := func() {
		if taskEventSub != nil {
			taskEventSub.Unsubscribe()
		}
		logger.Info().Msg("subscribing to TaskSubmitted events...")
		// Get current block number to start watching from
		subCtx, cancel := context.WithTimeout(appContext, 15*time.Second)
		defer cancel()
		blockNo, err := rpcClient.Client.BlockNumber(subCtx)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to get current block number for event subscription")
		}

		taskEventSub, err = engineContract.WatchTaskSubmitted(&bind.WatchOpts{
			Start:   &blockNo,
			Context: appContext,
		}, sinkTaskSubmitted, nil, nil, nil)
		if err != nil {
			logger.Error().Err(err).Msg("failed to subscribe to TaskSubmitted events")
		} else {
			logger.Info().Msg("subscribed to TaskSubmitted events")
		}
	}

	if cfg.ListenToTaskSubmitted {
		connectToTaskEvents()
	} else {
		logger.Warn().Msg("listening to TaskSubmitted events is disabled in config")
	}

	// Task Queue
	taskQueue, err := NewTaskQueue(logger, defaultMaxTasks, defaultTaskCacheSize) // Use constants or config
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create task queue")
	}
	logger.Info().Msg("task queue initialized")

	// --- Select and Start Mining Strategy ---
	var strategy MiningStrategy
	switch cfg.Strategies.Strategy {
	case "bulkmine":
		strategy = NewBulkMineStrategy(appQuit, appServices, miner, gpuPool, taskQueue)
	case "solutionsampler":
		// Need to connect to solution events *before* starting the strategy
		strategy = NewSolutionSamplerStrategy(appContext, appServices, miner, gpuPool, taskQueue)
	case "task": //  "task" is the simple listen-and-mine strategy
		if !cfg.ListenToTaskSubmitted {
			logger.Fatal().Msg("strategy 'task' requires 'ListenToTaskSubmitted' to be enabled in config")
		}
		// Simple strategy: just process tasks from the event stream directly
		logger.Warn().Msg("strategy 'task' currently behaves like 'bulkmine' fed by events")
		strategy = NewBulkMineStrategy(appContext, appServices, miner, gpuPool, taskQueue)

	default:
		logger.Fatal().Str("strategy", cfg.Strategies.Strategy).Msg("unknown or unsupported mining strategy specified in config")
	}

	// Start the selected strategy's workers
	err = strategy.Start()
	if err != nil {
		logger.Fatal().Err(err).Str("strategy", strategy.Name()).Msg("failed to start mining strategy")
	}

	dashboard := tui.NewDashboard()

	if !*headless {
		// this captures the log output and sends it to the logviewer
		logWriter, cleanup := tui.NewLogRouter()
		defer cleanup()

		// Update our log router to use the logviewer
		logWriter.SetView(dashboard.LogViewer.CustomTextView)
		go func() {
			dashboard.Run()
			// disable writing to log viwer on exit
			// we don't set the view to nil as this is unsafe instead we use an atomic bool
			logWriter.Headless.Store(true)
			appCancel()
		}()

		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					// GPU Metrics
					gpuMetrics := gpuPool.GetGPUInfoForMetrics() // Use GPUPool method
					dashboard.Updates <- tui.StateUpdate{Type: tui.UpdateGPUs, Payload: gpuMetrics}

					dashboard.Updates <- tui.StateUpdate{
						Type:    tui.UpdateGPUs,
						Payload: gpuMetrics,
					}

					// Update validator metrics
					validatorMetrics := tui.ValidatorMetrics{
						SessionTime:      manager.GetSessionTime(),      // Call new manager method
						SolvedLastMinute: manager.GetSolvedLastMinute(), // Call new manager method
						SolutionsLastMinute: struct {
							Success int64
							Total   int64
							Rate    float64
						}{
							Success: manager.GetSuccessCount(), // Call new manager method
							Total:   manager.GetTotalCount(),   // Call new manager method
							Rate:    manager.GetSuccessRate(),  // Call new manager method
						},
						AverageSolutionRate:    manager.GetAverageSolutionRate(),    // Call new manager method
						AverageSolutionsPerMin: manager.GetAverageSolutionsPerMin(), // Call new manager method
						AverageSolvesPerMin:    manager.GetAverageSolvesPerMin(),    // Call new manager method
					}
					dashboard.Updates <- tui.StateUpdate{
						Type:    tui.UpdateValidatorMetrics,
						Payload: validatorMetrics,
					}

					// // Update financial metrics
					// dashboard.updates <- StateUpdate{
					// 	Type: UpdateFinancialMetrics,
					// 	Payload: FinancialMetrics{
					// 		TokenIncomePerMin:  float64(count) * 0.1,
					// 		TokenIncomePerHour: float64(count) * 6.0,
					// 		TokenIncomePerDay:  float64(count) * 144.0,
					// 		IncomePerMin:       float64(count) * 0.5,
					// 		IncomePerHour:      float64(count) * 30.0,
					// 		IncomePerDay:       float64(count) * 720.0,
					// 		ProfitPerMin:       float64(count) * 0.3,
					// 		ProfitPerHour:      float64(count) * 18.0,
					// 		ProfitPerDay:       float64(count) * 432.0,
					// 	},
					// }
				}
			}
		}()

	} else {
		logger.Info().Msg("running in headless mode; dashboard disabled")
	}

	maxHeaderBackoff := 30 * time.Second
	currentHeaderBackoff := 1 * time.Second
	maxTaskEventBackoff := 30 * time.Second
	currentTaskEventBackoff := 1 * time.Second

	if cfg.ListenToTaskSubmitted {
		appQuitWG.Add(1)
		go func() {
			defer appQuitWG.Done()

			for {
				select {
				case <-appQuit.Done():
					logger.Info().Msg("task event subscription closed, exiting")
					return
				case event := <-sinkTaskSubmitted:
					if event == nil {
						continue
					} // Skip if channel closed or nil event sent

					currentTaskEventBackoff = 1 * time.Second // Reset backoff
					taskId := task.TaskId(event.Id)
					ts := &TaskSubmitted{
						TaskId: taskId,
						TxHash: event.Raw.TxHash,
					}
					logger.Info().Str("taskid", taskId.String()).Str("txHash", event.Raw.TxHash.Hex()).Uint64("block", event.Raw.BlockNumber).Msg("received TaskSubmitted event")

					// Add to task queue (AddTask handles deduplication)
					taskQueue.AddTask(ts)

				case err := <-taskEventSub.Err():
					if err == nil {
						continue
					}
					logger.Warn().Err(err).Msgf("task event subscription error, retrying connection in %s", currentTaskEventBackoff)
					time.Sleep(currentTaskEventBackoff)
					currentTaskEventBackoff = (currentTaskEventBackoff * 2) + time.Duration(rand.Intn(500))*time.Millisecond
					if currentTaskEventBackoff > maxTaskEventBackoff {
						currentTaskEventBackoff = maxTaskEventBackoff
					}
					connectToTaskEvents()
				}
			}
		}()
	}

	for {

		select {

		case <-appQuit.Done():
			logger.Info().Msg("shutting down main loop")
			goto exit_app
		case err := <-newHeadSub.Err():
			if err == nil {
				continue
			}
			logger.Warn().Msgf("new head sub error: %v, will retry in %s", err, currentHeaderBackoff.String())
			newHeadSub.Unsubscribe()

			time.Sleep(currentHeaderBackoff)
			currentHeaderBackoff = (currentHeaderBackoff * 2) + time.Duration(rand.Intn(500))*time.Millisecond
			if currentHeaderBackoff > maxHeaderBackoff {
				currentHeaderBackoff = maxHeaderBackoff
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
			//currentBlockNumer = h.Number.Uint64()
		}
	}
exit_app:
	logger.Info().Msg("waiting for application workers to finish")

	// Stop the mining strategy (signals workers, waits for them)
	if strategy != nil {
		strategy.Stop() // This should handle context cancellation and worker WaitGroup
	}

	// Wait for all workers to finish
	appQuitWG.Wait()

	logger.Info().Msg("bye! 👋")

	// for debugging purposes
	// Create a timeout channel to detect if the wait takes too long
	/* waitDone := make(chan struct{})
	go func() {
		appQuitWG.Wait()
		close(waitDone)
	}()

	// Wait for either completion or timeout
	select {
	case <-waitDone:
		logger.Info().Msg("all workers finished successfully")
	case <-time.After(10 * time.Second):
		logger.Warn().Msg("workers taking longer than expected to finish, dumping goroutine stacks for debugging")

		numGoroutines := runtime.NumGoroutine()
		logger.Warn().Int("count", numGoroutines).Msg("number of active goroutines")

		var buf bytes.Buffer

		pprof.Lookup("goroutine").WriteTo(&buf, 1)

		logger.Warn().Msgf("goroutine stacks:\n%s", buf.String())

		// Continue with the wait
		logger.Warn().Msg("continuing to wait for goroutines to finish")
	} */
}
