package main

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"errors"
	"flag"
	"fmt"
	cmn "gobius/common"
	"runtime"
	"runtime/pprof"

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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
// type Miner struct {
// 	services         *Services
// 	validator        IValidator
// 	engineAbi        *abi.ABI
// 	submitMethod     *abi.Method
// 	bulkSubmitMethod *abi.Method
// }

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

	// this captures the log output and sends it to the logviewer
	logWriter, cleanup := tui.NewLogRouter()
	defer cleanup()

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

	appContext, appServices, err := NewApplicationContext(rpcClient, sqlite, logger, cfg, ipfsOracle, context.Background(), appQuit)

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

		logWriter.RestoreOutputs()

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
			// get remove mode from args
			removeMode := false
			if len(args) > 1 {
				filename := args[1]
				if len(args) > 2 {
					removeMode, err = strconv.ParseBool(args[2])
					if err != nil {
						log.Fatalf("invalid remove mode value: %v", err)
					}
				}
				importUnsolvedTasks(appQuit, filename, removeMode, &logger, appContext)
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
		case "exportunsolved":
			var startBlock, endBlock int64
			var err error
			var senderFilter common.Address

			if len(args) < 2 {
				log.Fatal("unsolvedtasks requires at least a startblock argument")
			}

			startBlock, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				log.Fatalf("Invalid startblock value: %v", err)
			}

			if len(args) >= 3 {
				endBlock, err = strconv.ParseInt(args[2], 10, 64)
				if err != nil {
					log.Fatalf("Invalid endblock value: %v", err)
				}
			} else {
				endBlock = 0 // Signal to getUnsolvedTasks to use the latest block
			}

			// Check for optional sender filter address
			if len(args) >= 4 {
				if !common.IsHexAddress(args[3]) {
					log.Fatalf("Invalid sender filter address: %s", args[3])
				}
				senderFilter = common.HexToAddress(args[3])
			} else {
				senderFilter = common.Address{} // Use zero address if not provided
			}
			// Define a reasonable initial block size for the scan
			initialBlockSize := int64(10000) // Example initial size
			// make unique filename
			filename := fmt.Sprintf("unsolvedtasks_%s.json", time.Now().Format("20060102150405"))
			exportUnsolvedTasks(appQuit, appServices, rpcClient, startBlock, endBlock, initialBlockSize, senderFilter, filename)
		case "verifyclaims":
			verifyClaims(&logger, appContext)
		case "taskcheck":
			taskCheck(&logger, appContext)
		case "verifyalltasks":
			// extract if user wants to run in dry mode
			dryRun := false
			if len(args) > 1 {
				dryRun, err = strconv.ParseBool(args[1])
				if err != nil {
					log.Fatalf("invalid dry run value: %v", err)
				}
			}
			verifyAllTasks(appContext, dryRun)
		case "verifysolutions":
			verifySolutions(appContext)
		case "verifycommitments":
			verifyCommitment(appContext)
		case "blockmonitor":
			blockMonitor(appContext, rpcClient)
		case "recoverstale":
			recoverStaleTasks(appContext)
		case "claimbatchinfo":
			getBatchPricingInfo(appContext)
		case "unsolvedtasks":
			var startBlock, endBlock int64
			var err error
			var senderFilter common.Address

			if len(args) < 2 {
				log.Fatal("unsolvedtasks requires at least a startblock argument")
			}

			startBlock, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				log.Fatalf("Invalid startblock value: %v", err)
			}

			if len(args) >= 3 {
				endBlock, err = strconv.ParseInt(args[2], 10, 64)
				if err != nil {
					log.Fatalf("Invalid endblock value: %v", err)
				}
			} else {
				endBlock = 0 // Signal to getUnsolvedTasks to use the latest block
			}

			// Check for optional sender filter address
			if len(args) >= 4 {
				if !common.IsHexAddress(args[3]) {
					log.Fatalf("Invalid sender filter address: %s", args[3])
				}
				senderFilter = common.HexToAddress(args[3])
			} else {
				senderFilter = common.Address{} // Use zero address if not provided
			}

			// Define a reasonable initial block size for the scan
			initialBlockSize := int64(10000) // Example initial size
			getUnsolvedTasks(appQuit, appServices, rpcClient, startBlock, endBlock, initialBlockSize, senderFilter)
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
				appServices.Validators.InitiateValidatorWithdraw(validator, amount)
			case "completewithdrawall":
				appServices.Validators.ValidatorWithdraw(validator)
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
			appServices.Validators.CancelValidatorWithdraw(validator, index)
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

			appServices.Validators.VoteOnContestation(validator, taskid, yea)
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
		case "autotasksubmit":
			var interval time.Duration
			if len(args) == 2 {
				var err error
				interval, err = time.ParseDuration(args[1])
				if err != nil {
					log.Fatalf("Invalid duration format for autotasksubmit: %v. Use format like '5s', '1m', '300ms'", err)
				}
			} else {
				log.Fatalf("autotasksubmit requires a duration argument (e.g., '10s')")
			}
			// Assuming RunAutoTaskSubmit is defined in commands.go (in the same package)
			RunAutoTaskSubmit(appQuit, appServices, interval)
		case "gas-stats":
			var fromBlock, endBlock int64
			if len(args) < 2 {
				log.Fatal("gas-stats requires at least a from-block argument")
			}

			fromBlock, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				log.Fatalf("Invalid from-block value: %v", err)
			}

			if len(args) >= 3 {
				endBlock, err = strconv.ParseInt(args[2], 10, 64)
				if err != nil {
					log.Fatalf("Invalid end-block value: %v", err)
				}
			} else {
				endBlock = 0 // Signal to calculateGasStats to use the latest block
			}

			err = calculateGasStats(appQuit, appServices, rpcClient, fromBlock, endBlock, logger)
			if err != nil {
				logger.Fatal().Err(err).Msg("error during gas stats calculation")
			}
		case "analyzereward":
			var fromBlock, endBlock, sampleRate int64
			var threshold float64
			var err error

			if len(args) < 5 {
				log.Fatal("analyzereward requires from-block, end-block, threshold, and sample-rate arguments")
			}

			fromBlock, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				log.Fatalf("Invalid from-block value: %v", err)
			}

			endBlock, err = strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				log.Fatalf("Invalid end-block value: %v", err)
			}

			threshold, err = strconv.ParseFloat(args[3], 64)
			if err != nil {
				log.Fatalf("Invalid threshold value: %v", err)
			}

			sampleRate, err = strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				log.Fatalf("Invalid sample-rate value: %v", err)
			}
			if sampleRate <= 0 {
				log.Fatalf("sample-rate must be greater than 0")
			}

			err = analyzeRewardRecovery(appQuit, appServices, rpcClient, fromBlock, endBlock, threshold, sampleRate, logger)
			if err != nil {
				logger.Fatal().Err(err).Msg("error during reward recovery analysis")
			}
		default:
			log.Fatalf("unknown command: %s", command)
		}
		logger.Info().Msg("command executed successfully")
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

	modelAsBytes, _ := cmn.ConvertTaskIdString2Bytes(cfg.Strategies.Model)
	totalReward, err := appServices.Engine.GetModelReward(modelAsBytes)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get model reward")
	}

	rewardInAIUS := appServices.Config.BaseConfig.BaseToken.ToFloat(totalReward)
	logger.Info().Str("model", cfg.Strategies.Model).Str("reward", fmt.Sprintf("%.8g", rewardInAIUS)).Msg("selected strategy model reward")

	// miner, err := NewMinerEngine(appServices, manager, &appQuitWG)
	// if err != nil {
	// 	logger.Fatal().Err(err).Msg("could not create miner engine")
	// }

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
		if !appServices.Engine.VersionCheck(minerEngineVersion) {
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
	var connectErr error // Variable to capture connection error

	connectToHeaders := func() (ethereum.Subscription, error) {
		ctx, cancel := context.WithTimeout(appContext, 5*time.Second)
		defer cancel()

		// Attempt connection
		sub, err := rpcClient.Client.SubscribeNewHead(ctx, headers)
		if err != nil {
			logger.Error().Err(err).Msg("failed to subscribe to new headers")
			return nil, err // Return the error
		}
		logger.Info().Msg("subscribed to new headers")
		return sub, nil // Return the subscription and nil error
	}

	// Initial connection attempt
	newHeadSub, connectErr = connectToHeaders()

	// Select and start mining strategy
	var strategy MiningStrategy
	var strategyErr error
	switch cfg.Strategies.Strategy {
	case "bulkmine":
		strategy, strategyErr = NewBulkMineStrategy(appContext, appServices, manager, gpuPool)
	case "solutionsampler":
		// Need to connect to solution events *before* starting the strategy
		strategy, strategyErr = NewSolutionSamplerStrategy(appContext, appServices, manager, gpuPool)
	case "listen":
		strategy, strategyErr = NewListenStrategy(appContext, appServices, manager, gpuPool)
	case "automine":
		strategy, strategyErr = NewAutoMineStrategy(appContext, appServices, manager, gpuPool)
	default:
		// Use strategyErr to signal failure, consistent with other cases
		strategyErr = fmt.Errorf("unknown or unsupported mining strategy specified in config: %s", cfg.Strategies.Strategy)
	}

	// Check for errors during strategy initialization
	if strategyErr != nil {
		logger.Fatal().Err(strategyErr).Msg("failed to initialize mining strategy")
	}

	// Start the selected strategy
	err = strategy.Start()
	if err != nil {
		logger.Fatal().Err(err).Str("strategy", strategy.Name()).Msg("failed to start mining strategy")
	}

	dashboard := tui.NewDashboard()

	if !*headless {

		// Update our log router to use the logviewer
		logWriter.SetView(dashboard.LogViewer.CustomTextView)
		go func() {
			dashboard.Run()
			logWriter.StopTUIOutput()
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

	maxBackoff := 30 * time.Second    // Max backoff duration
	currentBackoff := 1 * time.Second // Initial backoff duration

	for {

		select {

		case <-appQuit.Done():
			logger.Info().Msg("shutting down main loop")
			goto exit_app
		case h := <-headers:
			if newHeadSub == nil {
				logger.Debug().Msg("header received but subscription inactive, skipping")
				continue // Skip processing if subscription is down
			}
			if h.BaseFee != nil {
				// update basefee
				rpcClient.SetBaseFee(h.BaseFee)
			}
		case err := <-func() <-chan error {
			if newHeadSub == nil { // Check if subscription is nil (e.g., initial connect failed)
				if connectErr != nil { // If there was an initial connect error, return a channel that immediately sends it
					errChan := make(chan error, 1)
					errChan <- connectErr
					close(errChan)
					connectErr = nil // Clear the initial error after sending
					return errChan
				}
				return nil // No active subscription and no pending initial error
			}
			return newHeadSub.Err() // Return the error channel of the active subscription
		}():
			if err == nil {
				// Channel closed unexpectedly? Treat as error.
				err = errors.New("header subscription error channel closed unexpectedly")
			}
			logger.Warn().Err(err).Msgf("header subscription error, attempting reconnect in %s", currentBackoff)

			// Cleanup existing subscription if it exists
			if newHeadSub != nil {
				newHeadSub.Unsubscribe()
				newHeadSub = nil
			}

			// Wait with backoff, checking for shutdown
			select {
			case <-time.After(currentBackoff):
				// Attempt reconnect
				newHeadSub, connectErr = connectToHeaders()
				// Adjust backoff for next potential failure
				currentBackoff = (currentBackoff * 2) + time.Duration(rand.Intn(500))*time.Millisecond
				if currentBackoff > maxBackoff {
					currentBackoff = maxBackoff
				}
				if newHeadSub != nil {
					logger.Info().Msg("reconnected to header subscription successfully")
					currentBackoff = 1 * time.Second // Reset backoff on success
				}
			case <-appQuit.Done():
				logger.Info().Msg("shutting down during header subscription reconnect backoff")
				goto exit_app
			}
		}
	}

exit_app:
	logger.Info().Msg("waiting for application workers to finish")

	// Wait for all workers to finish
	//appQuitWG.Wait()
	// for debugging purposes
	// Create a timeout channel to detect if the wait takes too long
	waitDone := make(chan struct{})
	go func() {
		// Stop the mining strategy (signals workers, waits for them)
		if strategy != nil {
			strategy.Stop()
		}
		appQuitWG.Wait()
		close(waitDone)
	}()

	// Wait for either completion or timeout
	select {
	case <-waitDone:
		logger.Info().Msg("all workers finished successfully")
	case <-time.After(60 * time.Second):
		logger.Warn().Msg("workers taking longer than expected to finish, dumping goroutine stacks for debugging")
		numGoroutines := runtime.NumGoroutine()
		logger.Warn().Int("count", numGoroutines).Msg("number of active goroutines")

		var buf bytes.Buffer
		pprof.Lookup("goroutine").WriteTo(&buf, 1)

		logger.Warn().Msgf("goroutine stacks:\n%s", buf.String())
	}

	logger.Info().Msg("bye! 👋")
}
