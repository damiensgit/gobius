package main

import (
	"bytes"
	"context"      // Add for sampling hash
	"database/sql" // Add for sampling hash conversion
	"errors"
	"fmt"
	"gobius/account"
	"gobius/bindings/engine"
	task "gobius/common" // Renamed import to avoid conflict
	"gobius/ipfs"
	"gobius/models"
	"gobius/utils"
	"math/rand"
	"sync"
	"time"

	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"

	"math/big"

	"github.com/ethereum/go-ethereum" // Import for ethereum.NotFound
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog"
)

// TODO: split this into multiple files

// TaskSubmitted represents a task identified by its ID and the transaction that submitted it.
// This struct is used across different producers and strategies.
type TaskSubmitted struct {
	TaskId [32]byte
	TxHash common.Hash
}

// goroutineRunner manages a WaitGroup for launching and waiting on goroutines.
type goroutineRunner struct {
	sync.WaitGroup
}

// Go starts a function in a new goroutine and handles WaitGroup accounting.
func (gr *goroutineRunner) Go(f func()) {
	gr.Add(1)
	go func() {
		defer gr.Done()
		f()
	}()
}

// MiningStrategy defines the common interface for different mining approaches.
type MiningStrategy interface {
	Start() error
	Stop()
	Name() string
}

// SolutionSubmitter is an interface for submitting solutions to the engine
type SolutionSubmitter interface {
	SignalCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) error
	SubmitIpfsCid(validator common.Address, taskId task.TaskId, cid []byte) error
	SubmitSolution(validator common.Address, taskId task.TaskId, cid []byte) error
}

// ErrTxDecodePermanent indicates a non-recoverable error during transaction decoding.
var ErrTxDecodePermanent = errors.New("permanent transaction decoding error")

const (
	defaultMaxSubscriptionBackoff     = 30 * time.Second
	defaultInitialSubscriptionBackoff = 1 * time.Second
)

// ConnectFunc is a function type responsible for establishing an event subscription.
// It should perform the actual contract `Watch...` call.
// The context passed to it has a timeout for the connection attempt.
type ConnectFunc func(ctx context.Context) (event.Subscription, error)

type subscriptionManager struct {
	logger      zerolog.Logger
	parentCtx   context.Context    // The context from the owner (e.g., producer)
	ctx         context.Context    // Internal context for the manager's loop
	cancel      context.CancelFunc // Cancels the internal context
	wg          *sync.WaitGroup    // Pointer to the owner's WaitGroup
	connectFunc ConnectFunc        // The function to establish the subscription
	eventName   string             // Name for logging (e.g., "TaskSubmitted")

	maxBackoff     time.Duration
	initialBackoff time.Duration

	mu           sync.Mutex
	subscription event.Subscription // Current active subscription

	stopOnce sync.Once
	goroutineRunner
}

// NewSubscriptionManager creates a new subscription manager.
func NewSubscriptionManager(parentCtx context.Context, wg *sync.WaitGroup, log zerolog.Logger, eventName string, connectFunc ConnectFunc) *subscriptionManager {
	ctx, cancel := context.WithCancel(parentCtx)
	return &subscriptionManager{
		logger:          log.With().Str("component", "subscriptionmanager").Str("event", eventName).Logger(),
		parentCtx:       parentCtx,
		ctx:             ctx,
		cancel:          cancel,
		wg:              wg,
		connectFunc:     connectFunc,
		eventName:       eventName,
		maxBackoff:      defaultMaxSubscriptionBackoff,
		initialBackoff:  defaultInitialSubscriptionBackoff,
		goroutineRunner: goroutineRunner{},
	}
}

// Start launches the background goroutine to manage the subscription.
func (sm *subscriptionManager) Start() {
	sm.logger.Info().Msg("starting")
	sm.Go(sm.manageLoop)
}

// Stop signals the manager loop to terminate and unsubscribes if needed.
func (sm *subscriptionManager) Stop() {
	sm.stopOnce.Do(func() {
		sm.logger.Info().Msg("stopping")
		// Cancel the internal context to signal the loop
		if sm.cancel != nil {
			sm.cancel()
		}
		sm.mu.Lock()
		if sm.subscription != nil {
			sm.subscription.Unsubscribe()
			sm.subscription = nil
		}
		sm.mu.Unlock()

		// wait for the goroutine(s) to finish
		sm.Wait()

		sm.logger.Info().Msg("stopped")
	})
}

// manageLoop is the core goroutine managing the subscription lifecycle.
func (sm *subscriptionManager) manageLoop() {
	sm.logger.Info().Msg("management loop started")

	currentBackoff := sm.initialBackoff
	var errChan <-chan error

	// Initial connection attempt
	if !sm.connectWithTimeout() {
		sm.logger.Warn().Msg("initial connection failed, will retry in loop")
		// Allow loop to handle backoff before first real wait
	}

	for {
		sm.mu.Lock()
		actSub := sm.subscription // Get current subscription under lock
		sm.mu.Unlock()

		if actSub == nil {
			// Subscription is down, attempt reconnect after backoff
			sm.logger.Warn().Dur("wait", currentBackoff).Msg("subscription down, attempting reconnect")
			select {
			case <-time.After(currentBackoff):
				currentBackoff = (currentBackoff * 2) + time.Duration(rand.Intn(500))*time.Millisecond
				if currentBackoff > sm.maxBackoff {
					currentBackoff = sm.maxBackoff
				}
				if !sm.connectWithTimeout() {
					continue // Connection failed, retry after next backoff interval
				}
				// Connection succeeded, reset backoff and get error channel
				currentBackoff = sm.initialBackoff
				sm.mu.Lock()
				if sm.subscription != nil { // Check again in case Stop was called during connect
					errChan = sm.subscription.Err() // Assignment is now valid
				} else {
					errChan = nil // Should not happen if connect succeeded, but be safe
				}
				sm.mu.Unlock()

			case <-sm.ctx.Done():
				sm.logger.Info().Msg("shutting down during reconnect backoff")
				return
			}
		} else {
			// Subscription is presumed active, use its error channel
			errChan = actSub.Err() // Assignment is now valid
		}

		select {
		case <-sm.ctx.Done():
			sm.logger.Info().Msg("context cancelled, shutting down loop")
			sm.mu.Lock()
			if sm.subscription != nil {
				sm.subscription.Unsubscribe()
				sm.subscription = nil
			}
			sm.mu.Unlock()
			return

		case err, ok := <-errChan:
			if !ok {
				// Channel closed unexpectedly? Treat as an error.
				err = errors.New("subscription error channel closed unexpectedly")
			}
			if err != nil { // Could be nil if channel just closed, but check anyway
				sm.logger.Warn().Err(err).Msg("subscription error detected")
			}
			// Regardless of error content, the subscription is broken
			sm.mu.Lock()
			if sm.subscription != nil {
				sm.subscription.Unsubscribe()
				sm.subscription = nil
			}
			sm.mu.Unlock()
			// Loop will now enter the reconnect logic in the next iteration
		}
	}
}

// connectWithTimeout attempts to establish the subscription using the connectFunc.
func (sm *subscriptionManager) connectWithTimeout() bool {
	// Use the manager's internal context for the attempt
	// TODO: make this configurable
	connectCtx, cancel := context.WithTimeout(sm.ctx, 20*time.Second) // e.g., 20 sec timeout
	defer cancel()

	sm.logger.Info().Msg("attempting to connect subscription...")
	newSub, err := sm.connectFunc(connectCtx)

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			sm.logger.Warn().Err(err).Msg("connection attempt timed out or was cancelled")
		} else {
			sm.logger.Error().Err(err).Msg("connection function failed")
		}
		return false
	}

	if newSub == nil {
		sm.logger.Error().Msg("connection function returned nil subscription without error")
		return false
	}

	sm.logger.Info().Msg("connection successful")
	sm.mu.Lock()
	// Check if Stop was called while we were connecting
	if sm.ctx.Err() != nil {
		sm.logger.Warn().Msg("context cancelled during successful connection, unsubscribing immediately")
		sm.mu.Unlock()
		newSub.Unsubscribe()
		return false
	}
	sm.subscription = newSub
	sm.mu.Unlock()
	return true
}

// TaskProducer defines the interface for components that provide tasks to workers.
type TaskProducer interface {
	// GetTask attempts to retrieve the next available task.
	// It BLOCKS if no task is immediately available, returning only when:
	//   - A task is ready (returns (*TaskSubmitted, nil)).
	//   - The passed context is cancelled (returns (nil, context.Canceled or specific error)).
	//   - The producer is stopped permanently (returns (nil, specific error)).
	GetTask(ctx context.Context) (*TaskSubmitted, error)

	// Start initializes the producer (e.g., starts listening, prepares DB).
	Start(ctx context.Context) error

	// Stop cleanly shuts down the producer.
	Stop()

	// Name returns a descriptive name for the producer (for logging).
	Name() string
}

// Task handler function signature now includes its own context
type taskHandlerFunc func(workerId int, gpu *task.GPU, ts *TaskSubmitted, taskCtx context.Context)

// Base Strategy (has common worker logic)
type baseStrategy struct {
	ctx          context.Context // Strategy's operational context
	cancelFunc   context.CancelFunc
	services     *Services
	submitter    SolutionSubmitter
	gpuPool      *GPUPool
	logger       zerolog.Logger
	numWorkers   int
	stopOnce     sync.Once
	strategyName string
	// TODO: refactor this, a bit ugly to have it here and then pass it to the producer
	txParamCache *lru.Cache // Keep cache for decodeTransaction in base

	// Embed goroutineRunner for managing worker goroutines
	goroutineRunner

	// Fields for TaskSubmitted event tracking (for TxHash lookup)
	taskSubManager    *subscriptionManager             // For TaskSubmitted events
	taskSubmittedSink chan *engine.EngineTaskSubmitted // Sink for TaskSubmitted events

	// Fields for sampled verification (now part of base)
	verifySubManager  *subscriptionManager                 // For SolutionSubmitted events
	verifySink        chan *engine.EngineSolutionSubmitted // Sink for SolutionSubmitted events
	sampledVerifyChan chan *TaskSubmitted                  // Channel for sampled verification tasks to workers

	submitMethod     *abi.Method
	bulkSubmitMethod *abi.Method
}

func (b *baseStrategy) Go(f func()) {
	b.goroutineRunner.Go(f)
}

// newBaseStrategy initializes the common components for any mining strategy.
// Returns a pointer to baseStrategy to avoid copying lock values.
func newBaseStrategy(appCtx context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool, strategyName string) (*baseStrategy, error) {
	ctx, cancel := context.WithCancel(appCtx) // Create a derived context for the strategy
	numWorkers := gpuPool.NumGPUs() * services.Config.NumWorkersPerGPU

	// Initialize LRU cache for transaction parameters (kept in base)
	// Use a reasonable default size
	// TODO: make this configurable / incease based on numWorkers?
	const defaultCacheSize = 100_000 // track up to 100k tasks
	paramCache, err := lru.New(defaultCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize transaction parameter LRU cache: %w", err)
	}

	logger := services.Logger.With().Str("strategy", strategyName).Logger()

	// TODO: Add VerificationSampleRate to Config struct
	sampleRate := uint64(1000) // Default sample rate (1 in 1000)
	logger.Warn().Uint64("rate", sampleRate).Msg("Using default verification sample rate (1 in N) - TODO: Add VerificationSampleRate to config")
	// if services.Config.VerificationSampleRate > 0 { ... }

	if sampleRate == 0 {
		logger.Error().Msg("VerificationSampleRate cannot be zero, disabling sampling verification")
		// Set sampleRate extremely high to effectively disable? Or handle differently?
		sampleRate = ^uint64(0) // Max uint64
	}

	// Get the parsed ABI once at startup
	parsed, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get engine ABI: %w", err)
	}

	// Get the submitTask method once at startup
	submitMethod, exists := parsed.Methods["submitTask"]
	if !exists {
		return nil, fmt.Errorf("submitTask method not found in ABI")
	}

	// Get the bulkSubmitTask method once at startup
	bulkSubmitMethod, exists := parsed.Methods["bulkSubmitTask"]
	if !exists {
		return nil, fmt.Errorf("bulkSubmitTask method not found in ABI")
	}

	bs := baseStrategy{
		ctx:               ctx,
		cancelFunc:        cancel,
		services:          services,
		submitter:         submitter,
		gpuPool:           gpuPool,
		logger:            logger,
		numWorkers:        numWorkers,
		strategyName:      strategyName,
		txParamCache:      paramCache,
		goroutineRunner:   goroutineRunner{},
		taskSubmittedSink: make(chan *engine.EngineTaskSubmitted, numWorkers*2),     // Buffer sink
		verifySink:        make(chan *engine.EngineSolutionSubmitted, numWorkers*2), // Buffer sink
		sampledVerifyChan: make(chan *TaskSubmitted, numWorkers),                    // Buffer channel to workers
		submitMethod:      &submitMethod,
		bulkSubmitMethod:  &bulkSubmitMethod,
	}

	// Setup TaskSubmitted Subscription
	bs.taskSubManager = NewSubscriptionManager(bs.ctx, &bs.WaitGroup, bs.logger, "TaskSubmitted", func(connectCtx context.Context) (event.Subscription, error) {
		client := bs.services.OwnerAccount.Client.Client
		if client == nil {
			return nil, errors.New("ethereum client is nil for TaskSubmitted subscription")
		}
		engineContract, err := engine.NewEngine(bs.services.Config.BaseConfig.EngineAddress, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create engine contract instance for TaskSubmitted: %w", err)
		}
		blockNo, err := client.BlockNumber(connectCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to get block number for TaskSubmitted subscription: %w", err)
		}
		// WatchTaskSubmitted expects filters: id [][32]byte, model [][32]byte, sender []common.Address
		return engineContract.WatchTaskSubmitted(&bind.WatchOpts{Start: &blockNo, Context: bs.ctx}, bs.taskSubmittedSink, nil, nil, nil)
	})

	// Setup SolutionSubmitted (Verify) Subscription
	bs.verifySubManager = NewSubscriptionManager(bs.ctx, &bs.WaitGroup, bs.logger, "SolutionSubmitted", func(connectCtx context.Context) (event.Subscription, error) {
		client := bs.services.OwnerAccount.Client.Client
		if client == nil {
			return nil, errors.New("ethereum client is nil for SolutionSubmitted subscription")
		}
		engineContract, err := engine.NewEngine(bs.services.Config.BaseConfig.EngineAddress, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create engine contract instance for SolutionSubmitted: %w", err)
		}
		blockNo, err := client.BlockNumber(connectCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to get block number for SolutionSubmitted subscription: %w", err)
		}
		// WatchSolutionSubmitted expects filters: addr []common.Address, task [][32]byte
		// Pass nil for both to watch all events.
		return engineContract.WatchSolutionSubmitted(&bind.WatchOpts{Start: &blockNo, Context: bs.ctx}, bs.verifySink, nil, nil)
	})

	// Start subscriptions and processing loops
	bs.taskSubManager.Start() // Start TaskSubmitted listener
	bs.verifySubManager.Start()
	bs.Go(bs.processTaskSubmittedLoop) // Start loop to process TaskSubmitted events into cache
	bs.Go(func() { bs.processSampledSolutionsLoop(sampleRate) })

	return &bs, nil // Return pointer to bs
}

// gpuWorker is simplified to pull directly from the producer.
func (bs *baseStrategy) gpuWorker(workerId int, gpu *task.GPU, producer TaskProducer, taskHandler taskHandlerFunc) {
	workerLogger := bs.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("producer", producer.Name()).Logger()
	workerLogger.Info().Msg("started worker")

	// Ticker to periodically check the status of the GPU and re-enable it
	ticker := time.NewTicker(time.Minute * 5) // TODO: make this configurable
	defer ticker.Stop()

	for {
		// Check context before potentially blocking/sleeping
		select {
		case <-bs.ctx.Done():
			workerLogger.Info().Msg("shutting down worker (context cancelled)")
			return
		case <-ticker.C:
			if !gpu.IsEnabled() {
				workerLogger.Info().Msg("gpu disabled, re-enabling")
				gpu.ResetErrorState()
			}
		default: // Continue if context is active
		}

		if !gpu.IsEnabled() {
			workerLogger.Debug().Msg("gpu disabled, pausing worker")
			select {
			case <-time.After(5 * time.Second): // Check status periodically
				continue
			case <-bs.ctx.Done():
				workerLogger.Info().Msg("shutting down worker while paused (gpu disabled)")
				return
			}
		}

		workerLogger.Debug().Msg("requesting task from producer...")
		// Use the strategy's context (bs.ctx) to get the task, allowing cancellation here
		ts, err := producer.GetTask(bs.ctx)

		if err != nil {
			// Check if the error is due to context cancellation (worker should exit)
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				workerLogger.Info().Err(err).Msg("worker context cancelled while getting task, exiting")
				return // Exit worker
			}

			// Check if the error indicates the producer is permanently stopped (worker should exit)
			if errors.Is(err, ErrProducerStopped) {
				workerLogger.Info().Err(err).Msg("producer stopped, worker exiting")
				return // Exit worker
			}

			// Otherwise, assume transient error, log, wait, and retry
			workerLogger.Warn().Err(err).Msg("transient error getting task from producer, retrying after delay...")
			// Add a small delay to avoid tight looping on persistent transient errors
			select {
			case <-time.After(5 * time.Second):
				continue // Retry GetTask
			case <-bs.ctx.Done():
				workerLogger.Info().Msg("worker context cancelled during retry delay, exiting")
				return // Exit if cancelled during sleep
			}
		}

		// If GetTask returns without error, we have a valid task
		workerLogger.Debug().Str("task", task.TaskId(ts.TaskId).String()).Msg("starting job")
		gpu.SetStatus("Mining")

		// Determine the context for the task handler based on config
		var taskCtx context.Context
		if bs.services.Config.Miner.WaitForTasksOnShutdown {
			// Use Background context to allow task completion on shutdown
			workerLogger.Debug().Msg("using background context for task execution (wait enabled)")
			taskCtx = context.Background()
		} else {
			// Use the worker's context to allow cancellation on shutdown
			workerLogger.Debug().Msg("using worker context for task execution (wait disabled)")
			taskCtx = bs.ctx
		}

		taskHandler(workerId, gpu, ts, taskCtx) // Pass the determined taskCtx

		workerLogger.Debug().Str("task", task.TaskId(ts.TaskId).String()).Msg("finished job processing")

		// Non-blocking check for SAMPLED verification tasks before getting next primary task
		select {
		case verifyTask, ok := <-bs.sampledVerifyChan:
			if !ok {
				workerLogger.Info().Msg("sampledVerifyChan closed, worker stopping verification checks")
				// Continue main loop, no more verification tasks will arrive.
			} else if verifyTask != nil {
				// Determine context for verification task
				var verifyTaskCtx context.Context
				// Use same WaitForTasksOnShutdown setting for verification?
				if bs.services.Config.Miner.WaitForTasksOnShutdown {
					verifyTaskCtx = context.Background()
				} else {
					verifyTaskCtx = bs.ctx // Allow cancellation if strategy stops
				}
				bs.handleVerificationTask(workerId, gpu, verifyTask, verifyTaskCtx)
			}
		default:
			// No verification task waiting, continue to next primary task
		}
	}
}

// start initializes workers and starts the producer.
func (bs *baseStrategy) start(producer TaskProducer, taskHandler taskHandlerFunc) error {
	bs.logger.Info().Int("workerspergpu", bs.services.Config.NumWorkersPerGPU).Int("gpus", bs.gpuPool.NumGPUs()).Msgf("starting %d workers for %s strategy", bs.numWorkers, bs.strategyName)

	// Start the producer itself
	err := producer.Start(bs.ctx) // Pass the strategy's context
	if err != nil {
		bs.logger.Error().Err(err).Msgf("failed to start producer %s", producer.Name())
		return err
	}
	bs.logger.Info().Msgf("producer %s started", producer.Name())

	gpus := bs.gpuPool.GetGPUs()
	for i := 0; i < bs.numWorkers; i++ {
		workerId := i
		gpuIndex := workerId / bs.services.Config.NumWorkersPerGPU
		if gpuIndex >= len(gpus) {
			bs.logger.Error().Int("workerId", workerId).Int("gpuIndex", gpuIndex).Int("numGpus", len(gpus)).Msg("logic error: worker assigned to non-existent GPU index")
			continue
		}
		gpu := gpus[gpuIndex]
		// Pass the producer to the worker
		bs.Go(func() {
			bs.gpuWorker(workerId, gpu, producer, taskHandler)
		})
	}

	return nil
}

// DecodeTaskTransaction decodes either a submitTask or bulkSubmitTask transaction
// and returns the relevant parameters.
func (bs *baseStrategy) decodeTaskTransaction(tx *types.Transaction) (*SubmitTaskParams, error) {
	// First verify this is a transaction to the Engine contract
	if tx.To() == nil || *tx.To() != bs.services.Config.BaseConfig.EngineAddress {
		return nil, fmt.Errorf("transaction not sent to Engine contract")
	}

	// Get the method signature (first 4 bytes)
	methodSig := tx.Data()[:4]
	var params []interface{}
	var err error

	// Verify method signature and unpack
	if bytes.Equal(bs.submitMethod.ID, methodSig) {
		params, err = bs.submitMethod.Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			return nil, fmt.Errorf("failed to unpack submitTask: %w", err)
		}
	} else if bytes.Equal(bs.bulkSubmitMethod.ID, methodSig) {
		params, err = bs.bulkSubmitMethod.Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			return nil, fmt.Errorf("failed to unpack bulkSubmitTask: %w", err)
		}
		// Note: bulkSubmitTask has an extra 'numTasks' param at the end, which we ignore here.
	} else {
		return nil, fmt.Errorf("transaction is not submitTask or bulkSubmitTask")
	}

	// Check if enough parameters were unpacked
	if len(params) < 5 {
		return nil, fmt.Errorf("unpacked fewer parameters than expected")
	}

	// Extract parameters (common to both methods up to this point)
	version, okV := params[0].(uint8)
	owner, okO := params[1].(common.Address)
	model, okM := params[2].([32]byte)
	fee, okF := params[3].(*big.Int)
	input, okI := params[4].([]byte)

	if !okV || !okO || !okM || !okF || !okI {
		return nil, fmt.Errorf("type assertion failed for one or more parameters")
	}

	submitTaskParams := &SubmitTaskParams{
		Version: version,
		Owner:   owner,
		Model:   model,
		Fee:     fee,
		Input:   input,
	}

	return submitTaskParams, nil
}

// DecodeSubmitTask decodes a submitTask transaction and returns the parameters
// This is a helper function to decode the parameters from a submitTask transaction to the engine contract only
// Deprecated: Use DecodeTaskTransaction instead which handles both single and bulk submit.
func (bs *baseStrategy) decodeSubmitTask(tx *types.Transaction, taskId task.TaskId) (*SubmitTaskParams, error) {
	// First verify this is a transaction to the Engine contract
	if tx.To() == nil || *tx.To() != bs.services.Config.BaseConfig.EngineAddress {
		return nil, fmt.Errorf("transaction not sent to Engine contract")
	}

	// Get the method signature (first 4 bytes)
	methodSig := tx.Data()[:4]

	// Verify this is a submitTask call by checking method signature
	if !bytes.Equal(bs.submitMethod.ID, methodSig) {
		return nil, fmt.Errorf("not a submitTask transaction")
	}

	// Now we can safely decode the parameters using cached method
	params, err := bs.submitMethod.Inputs.Unpack(tx.Data()[4:])
	if err != nil {
		return nil, err
	}

	// Now params contains your decoded parameters in order
	version := params[0].(uint8)
	owner := params[1].(common.Address)
	model := params[2].([32]byte)
	fee := params[3].(*big.Int)
	input := params[4].([]byte)

	submitTaskParams := &SubmitTaskParams{
		Version: version,
		Owner:   owner,
		Model:   model,
		Fee:     fee,
		Input:   input,
	}

	return submitTaskParams, nil
}

// decodeTransaction remains in baseStrategy as workers might need it via handleTask.
func (bs *baseStrategy) decodeTransaction(txHash common.Hash) (*SubmitTaskParams, error) {
	// 1. Check cache
	if cachedParams, found := bs.txParamCache.Get(txHash.String()); found {
		if params, ok := cachedParams.(*SubmitTaskParams); ok {
			bs.logger.Debug().Str("txHash", txHash.String()).Msg("using cached task parameters")
			return params, nil
		}
		bs.logger.Warn().Str("txHash", txHash.String()).Msg("invalid type found in txParamCache")
	}

	// 2. Fetch transaction if not cached
	txFetchStart := time.Now()
	// Use baseStrategy's context for the fetch
	tx, _, err := bs.services.OwnerAccount.Client.Client.TransactionByHash(bs.ctx, txHash)
	txFetchElapsed := time.Since(txFetchStart)
	if err != nil {
		// Check for specific, potentially permanent errors first
		if errors.Is(err, ethereum.NotFound) {
			bs.logger.Warn().Err(err).Str("txHash", txHash.String()).Msg("transaction not found")
			return nil, fmt.Errorf("transaction %s not found: %w", txHash.String(), err) // Wrap ethereum.NotFound
		}
		// Check for context errors (transient)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			bs.logger.Warn().Err(err).Str("txHash", txHash.String()).Msg("context cancelled or deadline exceeded while fetching transaction")
			return nil, fmt.Errorf("context error fetching transaction %s: %w", txHash.String(), err) // Wrap context errors
		}
		// Assume other errors are potentially transient network issues
		bs.logger.Error().Err(err).Str("txHash", txHash.String()).Msg("transient error fetching transaction")
		return nil, fmt.Errorf("transient error fetching transaction %s: %w", txHash.String(), err)
	}
	// Check if tx is nil even if err is nil (shouldn't happen with TransactionByHash, but defensive check)
	if tx == nil {
		bs.logger.Error().Str("txHash", txHash.String()).Msg("transaction not found (nil transaction returned)")
		// Treat as effectively not found, likely permanent issue with this hash on this node
		return nil, fmt.Errorf("transaction %s not found (nil tx): %w", txHash.String(), ethereum.NotFound)
	}

	bs.logger.Debug().Str("elapsed", txFetchElapsed.String()).Str("hash", txHash.String()).Msg("fetched transaction details")

	// 3. Decode transaction
	params, err := bs.decodeTaskTransaction(tx)
	if err != nil {
		bs.logger.Error().Err(err).Str("txHash", txHash.String()).Msg("permanent decode error for task transaction")
		// Wrap with ErrTxDecodePermanent
		return nil, fmt.Errorf("%w: %w", ErrTxDecodePermanent, err)
	}

	// 4. Store in cache
	bs.txParamCache.Add(txHash.String(), params)
	bs.logger.Debug().Str("txHash", txHash.String()).Msg("cached decoded task parameters")

	return params, nil
}

// SolveTask takes decoded task parameters and attempts to solve the task.
func (bs *baseStrategy) solveTask(ctx context.Context, taskId task.TaskId, params *SubmitTaskParams, gpu *task.GPU, validateOnly bool) ([]byte, error) {

	taskIdStr := taskId.String()

	inputRaw := string(params.Input) // Use input from params

	var result map[string]interface{}
	err := json.Unmarshal(params.Input, &result) // Use input from params

	if err != nil {
		bs.logger.Error().Err(err).Str("task", taskIdStr).Msg("could not unmarshal input from transaction parameters")
		return nil, err
	}

	bs.logger.Debug().Str("input", inputRaw).Str("task", taskIdStr).Msg("decoded task input")

	// taskInfo, err := m.services.Engine.LookupTask(taskId)
	// if err != nil {
	// 	m.services.Logger.Error().Err(err).Msg("could not lookup task")
	// 	return nil, err
	// }

	modelId := common.Bytes2Hex(params.Model[:]) // Use model from params
	model := models.ModelRegistry.GetModel(modelId)
	if model == nil {
		bs.logger.Error().Str("model", modelId).Str("task", taskIdStr).Msg("could not find model specified in task parameters")
		// Consider returning a specific error type here
		return nil, fmt.Errorf("model %s not found or enabled", modelId)
	}

	hydrated, err := model.HydrateInput(result, taskId.TaskId2Seed())

	if err != nil {
		bs.logger.Error().Err(err).Msg("could not hydrate input")
		return nil, err
	}

	output, err := json.Marshal(hydrated)
	if err != nil {
		bs.logger.Error().Err(err).Msg("could not marshal output")
		return nil, err
	}

	bs.logger.Debug().Str("output", string(output)).Msg("sending task to gpu")

	var cid []byte
	if bs.services.Config.EvilMode {
		cid, _ = hex.DecodeString("12206666666666666666666666666666666666666666666666666666666666666666")
		bs.logger.Warn().Str("cid", "0x"+hex.EncodeToString(cid)).Msg("evil mode enabled")
		duration := time.Duration(bs.services.Config.EvilModeMinTime+rand.Intn(bs.services.Config.EvilModeRandInt)) * time.Millisecond
		time.Sleep(duration)
	} else {
		//start := time.Now()
		if gpu.Mock {
			var data []byte
			data, err = gpu.GetMockCid(taskIdStr, hydrated)
			if err == nil {
				cid, err = ipfs.GetIPFSHashFast(data)
			}
		} else {
			cid, err = model.GetCID(ctx, gpu, taskIdStr, hydrated)
		}
		//elapsed := time.Since(start)
		//m.gpura.Add(elapsed)
		if err != nil {
			// Check if the error is due to context cancellation (timeout or explicit cancel)
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				// Log concise message for expected context cancellation
				bs.logger.Info().Str("task", taskIdStr).Msg("context cancelled or timed out during inference")
				return nil, err // Propagate the context error
			}

			// Handle other errors (GPU busy, actual inference errors, etc.)
			bs.logger.Error().Err(err).Str("task", taskIdStr).Msg("error on gpu during inference, incrementing error counter")
			gpu.IncrementErrorCount()
			return nil, err
		}
		//m.services.Logger.Debug().Str("cid", "0x"+hex.EncodeToString(cid)).Str("elapsed", elapsed.String()).Str("average", m.gpura.Average().String()).Msg("gpu finished & returned result")
	}

	if validateOnly {
		return cid, nil
	}

	validator := bs.services.Validators.GetNextValidatorAddress()

	commitment, err := utils.GenerateCommitment(validator, taskId, cid)
	if err != nil {
		bs.logger.Error().Err(err).Msg("error generating commitment hash")
		return nil, err
	}

	err = bs.submitter.SignalCommitment(validator, taskId, commitment)
	if err != nil {
		bs.logger.Error().Err(err).Msg("error signalling commitment to validator, skipping submitsolution")
		return nil, err
	}

	// we wont consider this a failure
	err = bs.submitter.SubmitIpfsCid(validator, taskId, cid)
	if err != nil {
		bs.logger.Warn().Err(err).Msg("ipfs cid submission failed")
	}

	// Use a separate goroutine without WaitGroup tracking for solution submission
	err = bs.submitter.SubmitSolution(validator, taskId, cid)
	if err != nil {
		bs.logger.Error().Err(err).Msg("solution submission failed")
	}

	return cid, err
}

// Stop stops the base strategy (cancels context, waits for workers).
// The Producer's Stop must be called separately by the strategy implementation.
func (bs *baseStrategy) Stop() {
	bs.stopOnce.Do(func() {
		bs.logger.Info().Msgf("stopping %s strategy workers", bs.strategyName)
		if bs.cancelFunc != nil {
			bs.cancelFunc()
		}
		bs.Wait() // Wait for worker goroutines using embedded WaitGroup
		bs.logger.Info().Msgf("all %s strategy workers stopped", bs.strategyName)
	})
}

// StorageProducer polls storage and provides tasks via a buffered channel.
type StorageProducer struct {
	services *Services
	logger   zerolog.Logger
	ctx      context.Context    // Context managed by the strategy that creates it
	cancel   context.CancelFunc // Cancel function for the poller loop
	stopOnce sync.Once

	taskChan chan *TaskSubmitted // Buffered channel for workers
	poolSize int                 // Max size of the channel buffer

	// Embed goroutineRunner for managing the poller loop
	goroutineRunner
}

// ErrProducerStopped indicates that the TaskProducer has been permanently stopped
// and cannot provide more tasks.
var ErrProducerStopped = errors.New("producer stopped")

// NewStorageProducer creates a producer that polls the storage.
func NewStorageProducer(appCtx context.Context, services *Services, poolSize int) *StorageProducer {
	ctx, cancel := context.WithCancel(appCtx) // Create derived context for internal loop
	if poolSize <= 0 {
		poolSize = 10 // Default buffer size, matches worker count?
	}
	p := &StorageProducer{
		services:        services,
		logger:          services.Logger.With().Str("producer", "storage").Logger(),
		ctx:             ctx,
		cancel:          cancel,
		taskChan:        make(chan *TaskSubmitted, poolSize),
		poolSize:        poolSize,
		goroutineRunner: goroutineRunner{},
	}
	return p
}

func (p *StorageProducer) Name() string { return "storage" }

// Start begins the background storage polling loop.
func (p *StorageProducer) Start(ctx context.Context) error {
	p.logger.Info().Msg("starting")
	p.Go(p.storageQueuePollerLoop) // Use the Go wrapper
	return nil
}

// Stop signals the poller loop to stop and closes the task channel.
func (p *StorageProducer) Stop() {
	p.stopOnce.Do(func() {
		p.logger.Info().Msg("stopping")
		if p.cancel != nil {
			p.cancel() // Signal background loop to stop
		}
		p.Wait()          // Wait for poller loop using embedded WaitGroup
		close(p.taskChan) // Close channel after poller stops writing
		p.logger.Info().Msg("stopped")
	})
}

// storageQueuePollerLoop continuously polls the storage and fills the task channel.
func (p *StorageProducer) storageQueuePollerLoop() {
	p.logger.Info().Msg("storage queue poller loop started")

	emptyPollInterval := 1 * time.Second
	errorPollInterval := 5 * time.Second
	backpressurePollInterval := 100 * time.Millisecond // How often to check when channel is full

	for {
		select {
		case <-p.ctx.Done(): // Check for stop signal first
			p.logger.Info().Msg("storage queue poller loop stopping")
			return
		default:
			// Only poll the storages if the output channel has space.
			if len(p.taskChan) >= p.poolSize {
				p.logger.Debug().Int("chan_len", len(p.taskChan)).Int("chan_cap", p.poolSize).Dur("wait", backpressurePollInterval).Msg("task queue full, pausing queue poll")
				// Wait a short interval before checking again, respecting context
				select {
				case <-time.After(backpressurePollInterval):
					continue // Loop back to check context and channel length again
				case <-p.ctx.Done():
					p.logger.Info().Msg("storage queue poller loop stopping during backpressure wait")
					return
				}
			}

			// Channel has space, okay to try popping a task
			p.logger.Debug().Msg("popping task from storage queue...")
			var taskId task.TaskId
			var txHash common.Hash
			var err error
			taskId, txHash, err = p.services.TaskStorage.PopTask(p.services.Config.PopTaskRandom)

			if err != nil {
				var sleepDuration time.Duration
				if errors.Is(err, sql.ErrNoRows) {
					sleepDuration = emptyPollInterval
					p.logger.Debug().Dur("wait", sleepDuration).Msg("storage queue empty, pausing poll")
				} else {
					sleepDuration = errorPollInterval
					p.logger.Error().Err(err).Dur("wait", sleepDuration).Msg("storage queue error, pausing poll")
				}

				// Wait before retrying, but check context during wait
				select {
				case <-time.After(sleepDuration):
					continue // Retry polling
				case <-p.ctx.Done():
					p.logger.Info().Msg("storage queue poller loop stopping during sleep")
					return
				}
			}

			// Successfully popped a task
			ts := &TaskSubmitted{TaskId: taskId, TxHash: txHash}
			p.logger.Debug().Str("task", taskId.String()).Msg("popped task from storage, attempting to buffer")

			// Send to channel (blocks if full), but handle context cancellation
			select {
			case p.taskChan <- ts:
				p.logger.Debug().Str("task", taskId.String()).Int("chan_len", len(p.taskChan)).Msg("task buffered")
				// Immediately try to poll again if channel wasn't full
			case <-p.ctx.Done():
				p.logger.Warn().Str("task", taskId.String()).Msg("storage queue poller stopping, discarding popped task")
				return // Exit loop
			}
		}
	}
}

// GetTask waits for a task from the internal channel or context cancellation.
func (p *StorageProducer) GetTask(ctx context.Context) (*TaskSubmitted, error) {
	p.logger.Debug().Int("chan_len", len(p.taskChan)).Msg("worker requesting task")
	select {
	case <-p.ctx.Done(): // Producer context stopping
		return nil, ErrProducerStopped // Use shared sentinel error
	case <-ctx.Done(): // Worker context stopping
		return nil, ctx.Err()
	case taskFromCh, ok := <-p.taskChan:
		if !ok {
			// Channel closed means producer is stopped
			return nil, ErrProducerStopped // Use shared sentinel error
		}
		p.logger.Info().Str("task", task.TaskId(taskFromCh.TaskId).String()).Int("chan_len", len(p.taskChan)).Msg("providing task from buffer")
		return taskFromCh, nil
	}
}

// EventProducer listens for on-chain TaskSubmitted events and provides them via a channel.
type EventProducer struct {
	services      *Services
	logger        zerolog.Logger
	ctx           context.Context    // Context managed by the strategy
	cancel        context.CancelFunc // Cancel function for the listener loop
	stopOnce      sync.Once
	taskChan      chan *TaskSubmitted // Buffered channel for workers
	sinkEvents    chan *engine.EngineTaskSubmitted
	maxBufferSize int

	// Subscription management
	subManager *subscriptionManager

	// Embed goroutineRunner for managing internal loops
	goroutineRunner
}

// NewEventProducer creates a producer that listens for on-chain events.
func NewEventProducer(appCtx context.Context, services *Services, bufferSize int) *EventProducer {
	ctx, cancel := context.WithCancel(appCtx)
	if bufferSize <= 0 {
		bufferSize = 100 // Default buffer size
	}
	p := &EventProducer{
		services:        services,
		logger:          services.Logger.With().Str("producer", "event").Logger(),
		ctx:             ctx,
		cancel:          cancel,
		taskChan:        make(chan *TaskSubmitted, bufferSize),              // Ensure taskChan is initialized
		sinkEvents:      make(chan *engine.EngineTaskSubmitted, bufferSize), // Buffer raw events too
		maxBufferSize:   bufferSize,
		goroutineRunner: goroutineRunner{}, // Initialize embedded runner
	}

	// Create the subscription manager, passing a closure for the connection logic
	connectFn := func(connectCtx context.Context) (event.Subscription, error) {
		client := p.services.OwnerAccount.Client.Client
		if client == nil {
			return nil, errors.New("ethereum client is nil")
		}
		engineAddr := p.services.Config.BaseConfig.EngineAddress
		engineContract, err := engine.NewEngine(engineAddr, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create engine contract instance: %w", err)
		}
		// Get block number within the connect timeout context
		blockNo, err := client.BlockNumber(connectCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to get block number for subscription: %w", err)
		}
		// IMPORTANT: Use the manager's PARENT context (p.ctx) for the long-running Watch call
		// The connectCtx is only for the setup phase (getting blockNo etc.)
		return engineContract.WatchTaskSubmitted(&bind.WatchOpts{
			Start:   &blockNo,
			Context: p.ctx, // Use producer's main context for watch duration
		}, p.sinkEvents, nil, nil, nil) // Use the WaitGroup from the embedded runner
	}

	p.subManager = NewSubscriptionManager(p.ctx, &p.WaitGroup, p.logger, "TaskSubmitted", connectFn)

	return p
}

func (p *EventProducer) Name() string { return "event" }

// Start begins the event listener loop.
func (p *EventProducer) Start(ctx context.Context) error {
	// Start the subscription manager (manages connection/reconnection)
	p.subManager.Start()
	// Start the loop to process events received in the sink
	p.Go(p.processEventsLoop) // Use Go wrapper
	return nil
}

// Stop signals the listener loop to stop and closes the task channel.
func (p *EventProducer) Stop() {
	p.stopOnce.Do(func() {
		p.logger.Info().Msg("stopping")
		// Stop the subscription manager first (cancels context, unsubscribes)
		p.subManager.Stop()
		if p.cancel != nil {
			p.cancel() // Signal listener loop
		}
		p.Wait()          // Wait for manager AND processEventsLoop
		close(p.taskChan) // Close channel after all goroutines exit
		p.logger.Info().Msg("stopped")
	})
}

// GetTask waits for a task from the internal channel or context cancellation.
func (p *EventProducer) GetTask(ctx context.Context) (*TaskSubmitted, error) {
	p.logger.Debug().Int("chan_len", len(p.taskChan)).Msg("worker requesting task")
	select {
	case <-p.ctx.Done(): // Producer context stopping
		return nil, ErrProducerStopped // Use shared sentinel error
	case <-ctx.Done(): // Worker context stopping
		return nil, ctx.Err()
	case taskFromCh, ok := <-p.taskChan:
		if !ok {
			// Channel closed means producer is stopped
			return nil, ErrProducerStopped // Use shared sentinel error
		}
		p.logger.Info().Str("task", task.TaskId(taskFromCh.TaskId).String()).Int("chan_len", len(p.taskChan)).Msg("providing task from event buffer")
		return taskFromCh, nil
	}
}

// processEventsLoop waits for events from the sink channel and pushes them to the task channel.
func (p *EventProducer) processEventsLoop() {
	p.logger.Info().Msg("starting event processing loop")

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info().Msg("shutting down event processing loop")
			return

		case event := <-p.sinkEvents:
			if event == nil {
				continue
			}

			ts := &TaskSubmitted{
				TaskId: event.Id,
				TxHash: event.Raw.TxHash,
			}
			taskIdStr := task.TaskId(ts.TaskId).String()

			p.logger.Info().Str("task", taskIdStr).Int("chan_len", len(p.taskChan)).Msg("received TaskSubmitted event")

			select {
			case p.taskChan <- ts:
				p.logger.Debug().Str("task", taskIdStr).Int("chan_len", len(p.taskChan)).Msg("event task buffered for worker")
			case <-p.ctx.Done():
				p.logger.Warn().Str("task", taskIdStr).Msg("context cancelled during event processing, discarding event")
				// Do not return here, allow loop to continue checking context
			default:
				p.logger.Warn().Str("task", taskIdStr).Int("buffer_size", p.maxBufferSize).Msg("event task channel full, discarding new event")
			}
		}
	}
}

// BulkMine Strategy: Uses StorageProducer.
type BulkMineStrategy struct {
	*baseStrategy
	producer TaskProducer
}

// NewBulkMineStrategy creates the strategy with a StorageProducer.
func NewBulkMineStrategy(appContext context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool) (*BulkMineStrategy, error) {
	base, err := newBaseStrategy(appContext, services, submitter, gpuPool, "bulkmine")
	if err != nil {
		return nil, err
	}
	// Size the producer buffer use numWorkers*2 for now?
	producer := NewStorageProducer(base.ctx, services, base.numWorkers*2)

	return &BulkMineStrategy{
		baseStrategy: base, // Use the pointer directly
		producer:     producer,
	}, nil
}

func (s *BulkMineStrategy) Name() string { return s.strategyName }

func (s *BulkMineStrategy) Start() error {
	s.logger.Info().Msgf("starting %s strategy", s.Name())
	return s.baseStrategy.start(s.producer, s.handleTask)
}

// handleTask processes a task received from the StorageProducer.
// Accepts taskCtx for the actual task processing.
func (s *BulkMineStrategy) handleTask(workerId int, gpu *task.GPU, ts *TaskSubmitted, taskCtx context.Context) {
	workerLogger := s.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", task.TaskId(ts.TaskId).String()).Logger()

	taskId := task.TaskId(ts.TaskId)

	requeueTask := func(taskId task.TaskId) {
		// Requeue ONLY because this strategy uses StorageProducer
		requeued, errDb := s.services.TaskStorage.RequeueTaskIfNoCommitmentOrSolution(taskId)
		if errDb != nil {
			workerLogger.Error().Err(errDb).Msg("failed to requeue task to storage")
		} else if requeued {
			workerLogger.Info().Msg("task requeued successfully to storage")
		} else {
			workerLogger.Warn().Msg("task not requeued (may have commitment/solution or other error)")
		}
	}

	// Decode the transaction (uses base cache)
	params, err := s.decodeTransaction(ts.TxHash)
	if err != nil {
		// Check if the error is permanent (decode issue or not found)
		if errors.Is(err, ErrTxDecodePermanent) || errors.Is(err, ethereum.NotFound) {
			workerLogger.Error().Err(err).Msg("permanent decode failure for task, dropping")
			// Do NOT requeue permanent decode failures
			// TODO: consider deleting the task from storage
		} else {
			// Assume other errors (like context cancelled or transient network issues) might be recoverable
			workerLogger.Warn().Err(err).Msg("transient error decoding transaction, requeueing task")
			requeueTask(taskId)
		}
		gpu.SetStatus("Idle") // Set GPU to error for ANY decode failure
		return
	}

	solveStart := time.Now()
	_, err = s.solveTask(taskCtx, taskId, params, gpu, false)
	solveElapsed := time.Since(solveStart)

	if err != nil {
		// Check if the error is specifically context cancellation or deadline exceeded
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			workerLogger.Info().Msg("task context cancelled, requeueing task")
			requeueTask(taskId)
			// Do not mark GPU as error, as it was context cancellation, not a GPU fault
			gpu.SetStatus("Idle") // Reset status as the task is being abandoned due to cancellation
		} else {
			// Handle other errors (genuine task processing failures)
			workerLogger.Error().Err(err).Msg("solve task failed, requeueing task")
			requeueTask(taskId)
			gpu.SetStatus("Error") // Mark GPU as having encountered an error
		}
	} else {
		workerLogger.Info().Str("elapsed", solveElapsed.String()).Msg("task solved successfully")
		s.gpuPool.AddSolveTime(solveElapsed)
		gpu.SetStatus("Idle")
	}
}

// Stop stops the base workers and the storage producer.
func (s *BulkMineStrategy) Stop() {
	s.logger.Info().Msgf("stopping %s strategy", s.Name())
	s.producer.Stop()     // Stop the producer first
	s.baseStrategy.Stop() // Then stop the base (workers)
	s.logger.Info().Msgf("%s strategy stopped", s.Name())
}

// AutoMine Strategy: Inherits from BulkMine, uses StorageProducer.
type AutoMineStrategy struct {
	*BulkMineStrategy // Embed BulkMineStrategy as a pointer
}

// NewAutoMineStrategy creates the strategy with a StorageProducer.
func NewAutoMineStrategy(appContext context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool) (*AutoMineStrategy, error) {
	if services.Config.Miner.BatchMode != 1 {
		return nil, errors.New("automine strategy requires batch mode to be enabled (solver.batchmode=1)")
	}
	if !services.Config.BatchTasks.Enabled {
		// this should be a warning
		services.Logger.Warn().Msg("automine strategy uses batch tasks, but batch tasks are not enabled (batchtasks.enabled=false)")
	}

	// TODO: validate that the automine model fee is set correctly - it needs to be the same as model.fee as per onchain ontract

	// Create the embedded BulkMineStrategy
	bulkStrategy, err := NewBulkMineStrategy(appContext, services, submitter, gpuPool)
	if err != nil {
		return nil, err
	}
	// Override the logger name
	bulkStrategy.logger = services.Logger.With().Str("strategy", "automine").Logger()
	bulkStrategy.strategyName = "automine" // Ensure base strategy name is also updated

	return &AutoMineStrategy{
		BulkMineStrategy: bulkStrategy,
	}, nil
}

// Listen Strategy uses EventProducer
type ListenStrategy struct {
	*baseStrategy
	producer TaskProducer
}

// NewListenStrategy creates the strategy with an EventProducer.
func NewListenStrategy(appContext context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool) (*ListenStrategy, error) {
	base, err := newBaseStrategy(appContext, services, submitter, gpuPool, "listen")
	if err != nil {
		return nil, err
	}
	// Size the producer buffer - needs config? Let's use numWorkers*2 for now.
	producer := NewEventProducer(base.ctx, services, base.numWorkers*2) // Use strategy's context

	return &ListenStrategy{
		baseStrategy: base, // Use the pointer directly
		producer:     producer,
	}, nil
}

func (s *ListenStrategy) Name() string { return s.strategyName }

// Start starts the base workers and the event producer.
func (s *ListenStrategy) Start() error {
	s.logger.Info().Msgf("starting %s strategy", s.Name())
	return s.baseStrategy.start(s.producer, s.handleTask)
}

// handleTask processes a task received from the EventProducer.
// Accepts taskCtx for the actual task processing.
func (s *ListenStrategy) handleTask(workerId int, gpu *task.GPU, ts *TaskSubmitted, taskCtx context.Context) {
	workerLogger := s.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", task.TaskId(ts.TaskId).String()).Logger()

	params, err := s.decodeTransaction(ts.TxHash)
	if err != nil {
		workerLogger.Error().Err(err).Msg("decode transaction failed (event task)")
		gpu.SetStatus("Error")
		return
	}

	solveStart := time.Now()
	_, err = s.solveTask(taskCtx, ts.TaskId, params, gpu, false)
	solveElapsed := time.Since(solveStart)

	if err != nil {
		workerLogger.Error().Err(err).Msg("solve task failed (event task - not requeued)")
		// Event tasks are ephemeral, not requeued
		gpu.SetStatus("Error")
	} else {
		workerLogger.Info().Str("elapsed", solveElapsed.String()).Msg("task solved successfully (event task)")
		s.gpuPool.AddSolveTime(solveElapsed)
		gpu.SetStatus("Idle")
	}
}

// Stop stops the base workers and the event producer.
func (s *ListenStrategy) Stop() {
	s.logger.Info().Msgf("stopping %s strategy", s.Name())
	s.producer.Stop()     // Stop producer
	s.baseStrategy.Stop() // Stop base (workers)
	s.logger.Info().Msgf("%s strategy stopped", s.Name())
}

// SolutionEventProducer listens for SolutionSubmitted events, looks up TxHash, and produces validation tasks.
// It collects samples over a period and dispatches them in batches.
type SolutionEventProducer struct {
	services     *Services
	logger       zerolog.Logger
	txParamCache *lru.Cache // Reference to baseStrategy's cache
	ctx          context.Context
	cancel       context.CancelFunc
	stopOnce     sync.Once
	taskChan     chan *TaskSubmitted // Output channel for validation tasks
	bufferSize   int

	// Event handling specific fields
	sinkTaskSubmitted     chan *engine.EngineTaskSubmitted     // Sink for TaskSubmitted events
	sinkSolutionSubmitted chan *engine.EngineSolutionSubmitted // Sink for SolutionSubmitted events
	maxTaskSampleSize     int                                  // Max samples to keep per period

	// Subscription management
	taskSubManager     *subscriptionManager // Manages TaskSubmitted subscription
	solutionSubManager *subscriptionManager // Manages SolutionSubmitted subscription

	// Sample collection state
	sampleMutex    sync.Mutex
	tasksSamples   []*TaskSubmitted // Using pointer to avoid copying large structs
	sampleIndex    int
	dispatchTicker *time.Ticker

	// Embed goroutineRunner for managing internal loops
	goroutineRunner
}

// NewSolutionEventProducer creates a producer for the SolutionSampler strategy.
// It requires access to the baseStrategy's txParamCache.
func NewSolutionEventProducer(appCtx context.Context, services *Services, txParamCache *lru.Cache, bufferSize, sampleSize int, dispatchInterval time.Duration) *SolutionEventProducer {
	ctx, cancel := context.WithCancel(appCtx)
	if bufferSize <= 0 {
		bufferSize = 100 // Default buffer size
	}
	if sampleSize <= 0 {
		sampleSize = 50 // Default sample size
	}
	if dispatchInterval <= 0 {
		dispatchInterval = 1 * time.Minute // Default dispatch interval
	}
	if txParamCache == nil {
		services.Logger.Error().Msg("SolutionEventProducer created without a valid txParamCache! It will not be able to produce tasks.")
	}

	p := &SolutionEventProducer{
		services:              services,
		logger:                services.Logger.With().Str("producer", "solutionevent").Logger(),
		txParamCache:          txParamCache,
		ctx:                   ctx,
		cancel:                cancel,
		taskChan:              make(chan *TaskSubmitted, bufferSize),
		bufferSize:            bufferSize,
		sinkTaskSubmitted:     make(chan *engine.EngineTaskSubmitted, 200),     // Sink channel for Task events
		sinkSolutionSubmitted: make(chan *engine.EngineSolutionSubmitted, 200), // Sink channel for Solution events
		maxTaskSampleSize:     sampleSize,
		tasksSamples:          make([]*TaskSubmitted, 0, sampleSize),
		dispatchTicker:        time.NewTicker(dispatchInterval),
		goroutineRunner:       goroutineRunner{}, // Initialize embedded runner
	}

	// ConnectFunc for TaskSubmitted events
	connectTaskFn := func(connectCtx context.Context) (event.Subscription, error) {
		client := p.services.OwnerAccount.Client.Client
		if client == nil {
			return nil, errors.New("ethereum client is nil")
		}
		engineAddr := p.services.Config.BaseConfig.EngineAddress
		engineContract, err := engine.NewEngine(engineAddr, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create engine contract instance: %w", err)
		}
		blockNo, err := client.BlockNumber(connectCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to get block number: %w", err)
		}
		return engineContract.WatchTaskSubmitted(&bind.WatchOpts{
			Start:   &blockNo,
			Context: p.ctx, // Use producer's main context for watch duration
		}, p.sinkTaskSubmitted, nil, nil, nil) // Use the WaitGroup from the embedded runner
	}
	p.taskSubManager = NewSubscriptionManager(p.ctx, &p.WaitGroup, p.logger, "TaskSubmitted", connectTaskFn)

	// ConnectFunc for SolutionSubmitted events
	connectSolutionFn := func(connectCtx context.Context) (event.Subscription, error) {
		client := p.services.OwnerAccount.Client.Client
		if client == nil {
			return nil, errors.New("ethereum client is nil")
		}
		engineAddr := p.services.Config.BaseConfig.EngineAddress
		engineContract, err := engine.NewEngine(engineAddr, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create engine contract instance: %w", err)
		}
		blockNo, err := client.BlockNumber(connectCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to get block number: %w", err)
		}
		return engineContract.WatchSolutionSubmitted(&bind.WatchOpts{
			Start:   &blockNo,
			Context: p.ctx, // Use producer's main context
		}, p.sinkSolutionSubmitted, nil, nil) // Use the WaitGroup from the embedded runner
	}
	p.solutionSubManager = NewSubscriptionManager(p.ctx, &p.WaitGroup, p.logger, "SolutionSubmitted", connectSolutionFn)

	return p
}

func (p *SolutionEventProducer) Name() string { return "solutionevent" }

// Start begins the subscription managers and the event processing/dispatching loops.
func (p *SolutionEventProducer) Start(ctx context.Context) error {
	if p.txParamCache == nil {
		p.logger.Error().Msg("cannot start: txParamCache is nil")
		return errors.New("SolutionEventProducer requires a valid txParamCache")
	}
	p.logger.Info().Msg("starting")

	p.taskSubManager.Start()     // Manages TaskSubmitted connection
	p.solutionSubManager.Start() // Manages SolutionSubmitted connection

	p.Go(p.processTaskSubmittedEvents)
	p.Go(p.processSolutionSubmittedEvents)
	p.Go(p.sampleDispatcherLoop)
	return nil
}

// Stop signals the loops to stop, stops the ticker/managers, and closes the output task channel.
func (p *SolutionEventProducer) Stop() {
	p.stopOnce.Do(func() {
		p.logger.Info().Msg("stopping")
		// Stop ticker first
		if p.dispatchTicker != nil {
			p.dispatchTicker.Stop()
		}
		// Stop subscription managers (unsubscribes, cancels internal context)
		if p.taskSubManager != nil {
			p.taskSubManager.Stop()
		}
		if p.solutionSubManager != nil {
			p.solutionSubManager.Stop()
		}
		// Cancel main producer context to signal processing loops
		if p.cancel != nil {
			p.cancel()
		}
		p.Wait()          // Waits for all goroutines started via p.Go AND the subscription managers
		close(p.taskChan) // Close channel after all goroutines exit
		p.logger.Info().Msg("stopped")
	})
}

// GetTask waits for a validation task generated from the dispatcher loop.
func (p *SolutionEventProducer) processTaskSubmittedEvents() {
	p.logger.Info().Msg("starting TaskSubmitted event processing loop (for caching)")

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info().Msg("shutting down TaskSubmitted event processing loop")
			return
		case event, ok := <-p.sinkTaskSubmitted:
			if !ok {
				p.logger.Info().Msg("TaskSubmitted sink channel closed, exiting loop")
				return
			}
			// should not happen but...
			if event == nil {
				continue
			}

			if p.txParamCache == nil {
				p.logger.Error().Msg("received TaskSubmitted event but txParamCache is nil, cannot cache")
				continue
			}
			p.logger.Debug().Str("task", task.TaskId(event.Id).String()).Msg("received TaskSubmitted, caching TxHash")
			p.txParamCache.Add(task.TaskId(event.Id).String(), event.Raw.TxHash)
		}
	}
}

// processSolutionSubmittedEvents waits for events from the sink and adds them to the sample pool.
func (p *SolutionEventProducer) processSolutionSubmittedEvents() {
	p.logger.Info().Msg("starting SolutionSubmitted event processing loop (for sampling)")

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info().Msg("shutting down SolutionSubmitted event processing loop")
			return
		case event, ok := <-p.sinkSolutionSubmitted:
			if !ok {
				p.logger.Info().Msg("SolutionSubmitted sink channel closed, exiting loop")
				return
			}
			if event == nil {
				continue
			}

			if p.txParamCache == nil {
				p.logger.Error().Msg("received SolutionSubmitted event but txParamCache is nil, cannot produce task")
				continue
			}

			taskIdStr := task.TaskId(event.Task).String()
			p.logger.Debug().Str("task", taskIdStr).Msg("received SolutionSubmitted event, looking up TxHash for sampling")

			// Look up TxHash from cache
			val, found := p.txParamCache.Get(taskIdStr)
			if !found {
				p.logger.Warn().Str("task", taskIdStr).Msg("solution event received but no TxHash found in cache, cannot validate/sample")
				continue
			}
			txHash, ok := val.(common.Hash)
			if !ok {
				p.logger.Error().Str("task", taskIdStr).Msgf("invalid type found in txParamCache for task, expected common.Hash, got %T", val)
				continue
			}

			ts := &TaskSubmitted{
				TaskId: event.Task,
				TxHash: txHash,
			}

			// Apply Reservoir Sampling
			p.sampleMutex.Lock()
			p.sampleIndex++
			if len(p.tasksSamples) < p.maxTaskSampleSize {
				p.tasksSamples = append(p.tasksSamples, ts)
				p.logger.Debug().Str("task", taskIdStr).Int("sample_size", len(p.tasksSamples)).Msg("added solution to sample pool")
			} else {
				j := rand.Intn(p.sampleIndex)
				if j < p.maxTaskSampleSize {
					p.tasksSamples[j] = ts
					p.logger.Debug().Str("task", taskIdStr).Int("replaced_index", j).Int("sample_size", len(p.tasksSamples)).Msg("replaced task in sample pool")
				}
			}
			p.sampleMutex.Unlock()
			// Note: Task is added to sample pool, dispatcher loop handles sending to workers.
		}
	}
}

// sampleDispatcherLoop periodically takes the collected samples and queues them for workers.
func (p *SolutionEventProducer) sampleDispatcherLoop() {
	p.logger.Info().Msg("starting sample dispatcher loop")

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info().Msg("stopping sample dispatcher loop")
			return
		case <-p.dispatchTicker.C:
			p.logger.Debug().Msg("dispatch ticker fired, processing samples")

			var newSamples []*TaskSubmitted

			// Lock, copy, clear internal sample buffer
			p.sampleMutex.Lock()
			if len(p.tasksSamples) > 0 {
				newSamples = make([]*TaskSubmitted, len(p.tasksSamples))
				copy(newSamples, p.tasksSamples)
				p.tasksSamples = p.tasksSamples[:0] // Clear while retaining capacity
				p.sampleIndex = 0
				p.logger.Info().Int("count", len(newSamples)).Msg("copied new samples for dispatch")
			} else {
				p.logger.Debug().Msg("no new samples collected in this period")
				// Even if no new samples, we might still need to trim the existing queue if it grew unexpectedly.
				// However, the main logic handles the combination and trimming naturally.
			}
			p.sampleMutex.Unlock()

			// --- Drain existing tasks from channel ---
			oldWaitingTasks := make([]*TaskSubmitted, 0, p.bufferSize) // Capacity hint
		drainingLoop:
			for {
				select {
				case task, ok := <-p.taskChan:
					if !ok {
						p.logger.Warn().Msg("task channel closed unexpectedly during drain")
						break drainingLoop // Exit if channel got closed
					}
					oldWaitingTasks = append(oldWaitingTasks, task)
				default: // Channel is empty
					break drainingLoop
				}
			}
			if len(oldWaitingTasks) > 0 {
				p.logger.Info().Int("count", len(oldWaitingTasks)).Msg("drained waiting tasks from channel")
			}

			// --- Combine new samples and old tasks (new samples first) ---
			combinedTasks := append(newSamples, oldWaitingTasks...)

			// --- Trim oldest tasks if combined list exceeds buffer size ---
			numToKeep := len(combinedTasks)
			droppedOldestCount := 0
			if numToKeep > p.bufferSize {
				droppedOldestCount = numToKeep - p.bufferSize
				combinedTasks = combinedTasks[:p.bufferSize] // Keep the first bufferSize elements (newest)
				numToKeep = p.bufferSize
				p.logger.Warn().Int("dropped_count", droppedOldestCount).Int("kept_count", numToKeep).Int("buffer_size", p.bufferSize).Msg("combined task list exceeded buffer size, oldest tasks dropped")
			}

			// --- Refill channel with the combined & trimmed list ---
			if numToKeep > 0 {
				p.logger.Info().Int("count", numToKeep).Msg("refilling channel with prioritized tasks...")
				for _, ts := range combinedTasks {
					if ts == nil {
						continue
					} // Safety check
					taskIdStr := task.TaskId(ts.TaskId).String()
					// Send should not block excessively as channel was drained and list trimmed to capacity,
					// but still respect context cancellation.
					select {
					case p.taskChan <- ts:
						p.logger.Debug().Str("task", taskIdStr).Int("q_len", len(p.taskChan)).Msg("added task back to channel")
					case <-p.ctx.Done():
						p.logger.Warn().Str("task", taskIdStr).Msg("context cancelled during channel refill, stopping refill")
						return // Stop refilling if context is cancelled
					}
				}
				p.logger.Info().Int("count", numToKeep).Int("final_q_len", len(p.taskChan)).Msg("finished refilling channel")
			} else if droppedOldestCount > 0 {
				// Log if we only dropped tasks and added nothing new (e.g., no new samples)
				p.logger.Info().Int("dropped_count", droppedOldestCount).Msg("only dropped oldest tasks, no new tasks to add")
			}
		}
	}
}

// GetTask waits for a validation task generated from the dispatcher loop.
func (p *SolutionEventProducer) GetTask(ctx context.Context) (*TaskSubmitted, error) {
	p.logger.Debug().Int("chan_len", len(p.taskChan)).Msg("worker requesting validation task")
	select {
	case <-p.ctx.Done(): // Producer context stopping
		return nil, ErrProducerStopped // Use shared sentinel error
	case <-ctx.Done(): // Worker context stopping
		return nil, ctx.Err()
	case taskFromCh, ok := <-p.taskChan:
		if !ok {
			// Channel closed means producer is stopped
			return nil, ErrProducerStopped // Use shared sentinel error
		}
		p.logger.Info().Str("task", task.TaskId(taskFromCh.TaskId).String()).Int("chan_len", len(p.taskChan)).Msg("providing validation task from event")
		return taskFromCh, nil
	}
}

// SolutionSampler Strategy listens for solutions submitted by others and validates them locally.
// Uses SolutionEventProducer.
type SolutionSamplerStrategy struct {
	*baseStrategy // Embed pointer
	producer      TaskProducer
}

// NewSolutionSamplerStrategy creates the strategy using SolutionEventProducer.
func NewSolutionSamplerStrategy(appCtx context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool) (*SolutionSamplerStrategy, error) {
	base, err := newBaseStrategy(appCtx, services, submitter, gpuPool, "solutionsampler")
	if err != nil {
		return nil, err
	}

	// Buffer size should match worker count or be slightly larger.
	bufferSize := base.numWorkers * 2
	sampleSize := bufferSize
	if sampleSize <= 0 {
		sampleSize = 50 // Default sample size if no workers
	}
	if bufferSize <= 0 { // If bufferSize (derived from workers) is 0, use default
		bufferSize = 20 // Default buffer
	}
	// TODO: Make dispatch interval configurable
	dispatchInterval := 1 * time.Minute

	producer := NewSolutionEventProducer(base.ctx, services, base.txParamCache, bufferSize, sampleSize, dispatchInterval)

	return &SolutionSamplerStrategy{
		baseStrategy: base, // Use the pointer directly
		producer:     producer,
	}, nil
}

func (s *SolutionSamplerStrategy) Name() string { return s.strategyName }

// Start initializes workers via baseStrategy, using the SolutionEventProducer.
func (s *SolutionSamplerStrategy) Start() error {
	s.logger.Info().Msgf("starting %s strategy", s.Name())
	// Use the standard base start method, passing the producer and the validation handler
	return s.baseStrategy.start(s.producer, s.handleValidationTask)
}

// handleValidationTask performs the validation logic.
// Accepts taskCtx for the actual task processing.
func (s *SolutionSamplerStrategy) handleValidationTask(workerId int, gpu *task.GPU, ts *TaskSubmitted, taskCtx context.Context) {
	workerLogger := s.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", task.TaskId(ts.TaskId).String()).Logger()
	workerLogger.Info().Msg("validating sampled task")
	// Set status immediately, might be changed on error
	gpu.SetStatus("Validating")

	// Decode the transaction (uses base cache or producer cache)
	params, err := s.decodeTransaction(ts.TxHash)
	if err != nil {
		// Check if the error is permanent
		if errors.Is(err, ErrTxDecodePermanent) || errors.Is(err, ethereum.NotFound) {
			workerLogger.Error().Err(err).Msg("permanent decode failure for solution task, skipping validation")
		} else {
			workerLogger.Warn().Err(err).Msg("transient error decoding solution task transaction, skipping validation")
		}
		gpu.SetStatus("Error") // Set status even for decode failure
		return
	}

	solveStart := time.Now()
	// Assign to ourCidBytes, handle error
	ourCidBytes, err := s.solveTask(taskCtx, ts.TaskId, params, gpu, true)
	solveElapsed := time.Since(solveStart)

	// Original error handling for SolveTask in validation (NO requeue)
	if err != nil {
		workerLogger.Error().Err(err).Msg("validation: solve task failed")
		if gpu.IsEnabled() {
			// Use more specific error status if possible
			// Check for context errors here too for more informative status?
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				gpu.SetStatus("Error")
			} else {
				gpu.SetStatus("Error")
			}
		}
		return // Exit validation on solve error
	}
	if ourCidBytes == nil {
		workerLogger.Error().Msg("validation: solve task did not return a CID")
		if gpu.IsEnabled() {
			gpu.SetStatus("Error")
		}
		return // Exit validation if no CID
	}

	workerLogger.Info().Str("elapsed", solveElapsed.String()).Msg("validation: task solved locally")

	// Fetch on-chain solution for comparison
	engineContract, err := engine.NewEngine(s.services.Config.BaseConfig.EngineAddress, s.services.OwnerAccount.Client.Client)
	if err != nil {
		workerLogger.Error().Err(err).Msg("validation: failed to create engine contract instance")
		gpu.SetStatus("Idle") // Reset status even if comparison fails
		return
	}

	callOpts := &bind.CallOpts{Context: s.ctx} // Use base strategy context
	res, err := engineContract.Solutions(callOpts, ts.TaskId)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			workerLogger.Warn().Err(err).Msg("validation: context cancelled during Solutions call")
		} else {
			workerLogger.Error().Err(err).Msg("validation: error getting on-chain solution info")
		}
		gpu.SetStatus("Idle") // Reset status
		return
	}

	if res.Blocktime == 0 {
		workerLogger.Warn().Msg("validation: no solution found on-chain (or call failed silently?), cannot compare")
		gpu.SetStatus("Idle") // Reset status
		return
	}

	solversCid := common.Bytes2Hex(res.Cid[:])
	ourCid := common.Bytes2Hex(ourCidBytes)

	workerLogger.Info().Str("our_cid", ourCid).Str("solver_cid", solversCid).Msg("comparing CIDs")

	if ourCid != solversCid {
		workerLogger.Warn().Msg("==================== CID MISMATCH DETECTED =====================")
		workerLogger.Warn().Msgf("  Task ID  : %s", task.TaskId(ts.TaskId).String())
		workerLogger.Warn().Msgf("  Our CID  : %s", ourCid)
		workerLogger.Warn().Msgf("  Their CID: %s", solversCid)
		workerLogger.Warn().Msgf("  Solver   : %s", res.Validator.String())
		workerLogger.Warn().Msgf("  Block    : %d", res.Blocktime)
		workerLogger.Warn().Msg("================================================================")
	} else {
		workerLogger.Info().Msg("validation: CID matches on-chain solution")
	}

	// Set status back to Idle after successful validation/comparison
	gpu.SetStatus("Idle")
}

// Stop stops the producer and the base workers.
func (s *SolutionSamplerStrategy) Stop() {
	s.logger.Info().Msgf("stopping %s strategy", s.Name())
	s.producer.Stop()     // Stop the specific producer first
	s.baseStrategy.Stop() // Then stop the base (waits for workers)
	s.logger.Info().Msgf("%s strategy stopped", s.Name())
}

// processSampledSolutionsLoop listens for SolutionSubmitted events, samples them based on a hash, and sends them for verification.
func (bs *baseStrategy) processSampledSolutionsLoop(sampleRate uint64) {
	bs.logger.Info().Uint64("sampleRate", sampleRate).Msg("starting sampled SolutionSubmitted processing loop")
	if sampleRate == 0 || sampleRate == ^uint64(0) { // Check if sampling is effectively disabled
		bs.logger.Warn().Msg("Sampling rate is zero or max uint64, verification sampling loop will not process events.")
		// Exit loop if disabled?
		// return
		// Or let it run but just log skip? Let's keep it running for now, but log verbosely.
	}

	for {
		select {
		case <-bs.ctx.Done():
			bs.logger.Info().Msg("shutting down sampled SolutionSubmitted processing loop")
			return
		case event, ok := <-bs.verifySink:
			if !ok {
				bs.logger.Info().Msg("SolutionSubmitted sink channel closed, exiting sampling loop")
				return
			}
			if event == nil {
				continue
			}

			taskIdStr := task.TaskId(event.Task).String()
			solverAddrStr := event.Addr.Hex()
			bs.logger.Debug().Str("task", taskIdStr).Str("solver", solverAddrStr).Msg("received SolutionSubmitted event")

			// Predictable Sampling Logic
			shouldSample := false
			if sampleRate > 0 && sampleRate != ^uint64(0) {
				hashData := append(event.Task[:], event.Addr.Bytes()...)
				hash := sha256.Sum256(hashData)
				// Use first 8 bytes of hash as uint64 for modulo operation
				sampleValue := binary.BigEndian.Uint64(hash[:8])
				if (sampleValue % sampleRate) == 0 {
					shouldSample = true
					bs.logger.Info().Str("task", taskIdStr).Str("solver", solverAddrStr).Uint64("hashValue", sampleValue).Msg("sampling solution for verification")
				}
			} else {
				bs.logger.Debug().Str("task", taskIdStr).Msg("Verification sampling disabled, skipping solution.")
				continue // Don't process if sampling disabled
			}

			if !shouldSample {
				bs.logger.Trace().Str("task", taskIdStr).Str("solver", solverAddrStr).Msg("solution not sampled")
				continue
			}

			// Solution IS sampled, proceed to get TxHash and queue for worker
			var txHash common.Hash
			found := false

			// --- TxHash Lookup: ONLY check txParamCache (Part 1) ---
			if bs.txParamCache != nil {
				if val, cacheFound := bs.txParamCache.Get(taskIdStr); cacheFound {
					if hashVal, ok := val.(common.Hash); ok {
						txHash = hashVal
						found = true
						bs.logger.Debug().Str("task", taskIdStr).Msg("found TxHash in txParamCache for sampled task")
					} else {
						bs.logger.Warn().Str("task", taskIdStr).Msgf("invalid type found in txParamCache for task, expected common.Hash, got %T", val)
					}
				}
			} else {
				bs.logger.Error().Str("task", taskIdStr).Msg("txParamCache is nil, cannot look up TxHash")
			}

			// --- DB Lookup REMOVED for Part 1 ---
			/*
				if !found { // If not found in cache, check DB (Part 2 logic)
					bs.logger.Debug().Str("task", taskIdStr).Msg("TxHash not found in cache, checking persistent storage...")
					// txHash, err := bs.services.YourStorageInterface.GetTxHashForTask(bs.ctx, event.Task)
					// if err == nil && txHash != (common.Hash{}) {
					// 	found = true
					// 	bs.logger.Info().Str("task", taskIdStr).Msg("found TxHash in DB for sampled task")
					// } else if err != nil && !errors.Is(err, YourStorageNotFoundError) {
					// 	bs.logger.Error().Err(err).Str("task", taskIdStr).Msg("error querying DB for TxHash")
					// }
				}
			*/

			if !found {
				bs.logger.Warn().Str("task", taskIdStr).Msg("TxHash not found in txParamCache for sampled task, skipping verification")
				continue
			}

			ts := &TaskSubmitted{
				TaskId: event.Task,
				TxHash: txHash,
			}

			// Send task to verify channel for workers (non-blocking)
			select {
			case bs.sampledVerifyChan <- ts:
				bs.logger.Info().Str("task", taskIdStr).Int("q_len", len(bs.sampledVerifyChan)).Msg("queued sampled task for verification")
			case <-bs.ctx.Done():
				bs.logger.Warn().Str("task", taskIdStr).Msg("context cancelled while queueing sampled verification task")
			default:
				// Log clearly that a sampled task is being dropped
				bs.logger.Error().Str("task", taskIdStr).Int("q_len", len(bs.sampledVerifyChan)).Int("q_cap", cap(bs.sampledVerifyChan)).Msg("SAMPLED verification task dropped, channel full!")
			}
		}
	}
}

// handleVerificationTask performs the validation logic for a sampled task.
func (bs *baseStrategy) handleVerificationTask(workerId int, gpu *task.GPU, ts *TaskSubmitted, taskCtx context.Context) {
	workerLogger := bs.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", task.TaskId(ts.TaskId).String()).Logger()
	workerLogger.Info().Msg("starting verification for sampled task")

	gpu.SetStatus("Verifying")
	defer gpu.SetStatus("Idle") // Ensure status is reset

	// Decode the transaction (uses base cache)
	params, err := bs.decodeTransaction(ts.TxHash)
	if err != nil {
		if errors.Is(err, ErrTxDecodePermanent) || errors.Is(err, ethereum.NotFound) {
			workerLogger.Error().Err(err).Msg("verification: permanent decode failure, skipping")
		} else {
			workerLogger.Warn().Err(err).Msg("verification: transient error decoding transaction, skipping")
		}
		return // Don't proceed with verification if decode fails
	}

	solveStart := time.Now()
	ourCidBytes, err := bs.solveTask(taskCtx, ts.TaskId, params, gpu, true) // true for validation mode
	solveElapsed := time.Since(solveStart)

	if err != nil {
		workerLogger.Error().Err(err).Msg("verification: local solve failed")
		if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) && gpu.IsEnabled() {
			gpu.SetStatus("Error") // Mark GPU error only if not cancelled and GPU enabled
		}
		return // Exit verification on solve error
	}
	if ourCidBytes == nil {
		workerLogger.Error().Msg("verification: local solve did not return a CID")
		if gpu.IsEnabled() {
			gpu.SetStatus("Error")
		}
		return // Exit validation if no CID
	}

	workerLogger.Info().Str("elapsed", solveElapsed.String()).Msg("verification: task solved locally")

	// Fetch on-chain solution for comparison
	engineContract, err := engine.NewEngine(bs.services.Config.BaseConfig.EngineAddress, bs.services.OwnerAccount.Client.Client)
	if err != nil {
		workerLogger.Error().Err(err).Msg("verification: failed to create engine contract instance")
		return
	}

	// Use taskCtx for RPC calls related to this specific verification
	callOpts := &bind.CallOpts{Context: taskCtx}
	res, err := engineContract.Solutions(callOpts, ts.TaskId)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			workerLogger.Warn().Err(err).Msg("verification: context cancelled during Solutions call")
		} else {
			workerLogger.Warn().Err(err).Msg("verification: error getting on-chain solution info (may not exist yet)")
		}
		return
	}

	if res.Blocktime == 0 {
		workerLogger.Warn().Msg("verification: no solution found on-chain (res.Blocktime is 0), cannot compare")
		return
	}

	solversCid := common.Bytes2Hex(res.Cid[:])
	ourCid := common.Bytes2Hex(ourCidBytes)

	workerLogger.Info().Str("our_cid", ourCid).Str("solver_cid", solversCid).Msg("comparing CIDs")

	// --- Validator Selection & Eligibility Check ---
	var selectedValidator *account.Account
	var bestStake = big.NewInt(0)
	eligibleFound := false

	for _, validator := range bs.services.Validators.validators {
		validatorAddr := validator.ValidatorAddress()
		workerLogger.Debug().Str("validator", validatorAddr.Hex()).Msg("Checking eligibility...")

		// Check if the validator can vote on this specific task
		votableStake, err := bs.services.Engine.Engine.ValidatorCanVote(callOpts, validatorAddr, ts.TaskId)
		if err != nil {
			workerLogger.Warn().Err(err).Str("validator", validatorAddr.Hex()).Msg("Failed to check ValidatorCanVote")
			continue // Skip this validator if check fails
		}

		if votableStake == nil || votableStake.Cmp(big.NewInt(0)) <= 0 {
			workerLogger.Debug().Str("validator", validatorAddr.Hex()).Msg("Validator not eligible or has insufficient stake to vote/contest this task")
			continue // Skip this validator
		}

		// Found an eligible validator, check if it has the best stake so far
		// For simplicity, using votableStake as the metric. Could also check total stake.
		if votableStake.Cmp(bestStake) > 0 {
			workerLogger.Debug().Str("validator", validatorAddr.Hex()).Str("votableStake", votableStake.String()).Msg("Found potentially better eligible validator")
			bestStake = votableStake
			selectedValidator = validator.Account
			eligibleFound = true
		}
	}

	if !eligibleFound {
		workerLogger.Warn().Msg("No eligible validator found with sufficient stake to vote/contest this task.")
		return // Cannot proceed
	}

	workerLogger.Info().Str("selectedValidator", selectedValidator.Address.Hex()).Str("effectiveStake", bestStake.String()).Msg("Selected best eligible validator for action")

	// --- Perform Action using Selected Validator ---
	if ourCid != solversCid {
		workerLogger.Warn().Msg("==================== CID MISMATCH DETECTED (Verification) ====================")
		workerLogger.Warn().Msgf("  Task ID  : %s", task.TaskId(ts.TaskId).String())
		workerLogger.Warn().Msgf("  Our CID  : %s", ourCid)
		workerLogger.Warn().Msgf("  Their CID: %s", solversCid)
		workerLogger.Warn().Msgf("  Solver   : %s", res.Validator.String()) // Validator field holds the solver addr
		workerLogger.Warn().Msgf("  Block    : %d", res.Blocktime)
		workerLogger.Warn().Msg("===========================================================================")

		// <<< CONTESTATION >>>
		workerLogger.Warn().Str("validator", selectedValidator.Address.Hex()).Msg("Initiating contestation...")
		// Use NonceManagerWrapper with the SELECTED validator account
		_, err = selectedValidator.NonceManagerWrapper(1, 425, 1.5, true, func(opts *bind.TransactOpts) (interface{}, error) {
			// Call SubmitContestation via the EngineTransactor
			return bs.services.Engine.Engine.EngineTransactor.SubmitContestation(opts, ts.TaskId)
		})
		if err != nil {
			workerLogger.Error().Err(err).Str("validator", selectedValidator.Address.Hex()).Msg("Failed to submit contestation transaction via NonceManagerWrapper")
			// txHash is not available here if the wrapper handles logging
		} else {
			// NonceManagerWrapper likely logs the hash, just log success indicator
			workerLogger.Info().Str("validator", selectedValidator.Address.Hex()).Msg("Contestation transaction submitted via NonceManagerWrapper (check wrapper logs for hash)")
		}
	} else {
		workerLogger.Info().Msg("verification: CID matches on-chain solution")

		// <<< VOTING >>>
		workerLogger.Info().Str("validator", selectedValidator.Address.Hex()).Msg("Initiating vote...")
		// Use NonceManagerWrapper with the SELECTED validator account
		_, err = selectedValidator.NonceManagerWrapper(1, 425, 1.5, true, func(opts *bind.TransactOpts) (interface{}, error) {
			// Call VoteOnContestation via the EngineTransactor
			return bs.services.Engine.Engine.EngineTransactor.VoteOnContestation(opts, ts.TaskId, true) // true for yea vote
		})
		if err != nil {
			workerLogger.Error().Err(err).Str("validator", selectedValidator.Address.Hex()).Msg("Failed to submit vote transaction via NonceManagerWrapper")
			// txHash is not available here
		} else {
			// NonceManagerWrapper likely logs the hash
			workerLogger.Info().Str("validator", selectedValidator.Address.Hex()).Msg("Vote (yea) transaction submitted via NonceManagerWrapper (check wrapper logs for hash)")
		}
	}
}

// processTaskSubmittedLoop listens to the task submission event sink and populates the txParamCache.
func (bs *baseStrategy) processTaskSubmittedLoop() {
	bs.logger.Info().Msg("starting TaskSubmitted processing loop (for cache)")
	for {
		select {
		case <-bs.ctx.Done():
			bs.logger.Info().Msg("shutting down TaskSubmitted processing loop")
			return
		case event, ok := <-bs.taskSubmittedSink:
			if !ok {
				bs.logger.Info().Msg("TaskSubmitted sink channel closed, exiting loop")
				return
			}
			if event == nil {
				continue
			}
			taskIdStr := task.TaskId(event.Id).String()
			bs.logger.Debug().Str("task", taskIdStr).Str("txHash", event.Raw.TxHash.Hex()).Msg("received TaskSubmitted event, adding/updating txParamCache")

			// Add/Update the cache. If the task ID was already present, this updates its TxHash
			// and makes it the most recently used entry.
			if bs.txParamCache != nil {
				bs.txParamCache.Add(taskIdStr, event.Raw.TxHash)
			}
		}
	}
}
