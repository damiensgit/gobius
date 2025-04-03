package main

import (
	"context"
	"database/sql"
	"fmt"
	"gobius/bindings/engine"
	task "gobius/common"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog"
)

// MiningStrategy defines the interface for different mining approaches.
type MiningStrategy interface {
	Start() error
	Stop()
	// Name returns the strategy name (e.g., "bulkmine", "solutionsampler")
	Name() string
	// HandleTaskSubmitted allows strategies to react to new tasks submitted on-chain.
	// Not all strategies need to implement this.
	HandleTaskSubmitted(event *engine.EngineTaskSubmitted)
}

type taskHandlerFunc func(workerId int, gpu *task.GPU, ts *TaskSubmitted)

// --- Base Strategy (Common worker logic) ---
type baseStrategy struct {
	ctx               context.Context // Strategy's operational context
	cancelFunc        context.CancelFunc
	services          *Services
	miner             *Miner
	gpuPool           *GPUPool
	taskQueue         *TaskQueue
	logger            zerolog.Logger
	numWorkers        int
	wg                sync.WaitGroup
	stopOnce          sync.Once
	strategyName      string
	txParamCache      *lru.Cache                       // Cache for TxHash -> *SubmitTaskParams (Using LRU)
	needsTaskEvents   bool                             // Flag indicating if the strategy needs TaskSubmitted events
	taskEventSub      event.Subscription               // Subscription object if listening
	sinkTaskSubmitted chan *engine.EngineTaskSubmitted // Channel for receiving events if listening
}

func (b *baseStrategy) Go(f func()) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		f()
	}()
}

// Generic implementation, can be overridden by specific strategies.
func (bs *baseStrategy) HandleTaskSubmitted(event *engine.EngineTaskSubmitted) {
	// Default: Do nothing or log a warning if unexpected
	// bs.logger.Debug().Str("task", task.TaskId(event.Id).String()).Msg("received TaskSubmitted event but strategy does not handle it")
}

func newBaseStrategy(appCtx context.Context, services *Services, miner *Miner, gpuPool *GPUPool, strategyName string, needsTaskEvents bool) (baseStrategy, error) {
	ctx, cancel := context.WithCancel(appCtx) // Create a derived context for the strategy
	numWorkers := gpuPool.NumGPUs() * services.Config.NumWorkersPerGPU

	var sink chan *engine.EngineTaskSubmitted
	if needsTaskEvents {
		sink = make(chan *engine.EngineTaskSubmitted, 1024)
	}

	// Initialize LRU cache for transaction parameters
	paramCacheSize := defaultTaskCacheSize // Use the constant defined in task_queue.go
	paramCache, err := lru.New(paramCacheSize)
	if err != nil {
		// Log fatal because cache is essential for performance
		// services.Logger.Fatal().Err(err).Msg("failed to initialize transaction parameter LRU cache")
		return baseStrategy{}, fmt.Errorf("failed to initialize transaction parameter LRU cache: %w", err) // Return error
	}

	// Initialize TaskQueue internally
	taskQueue, err := NewTaskQueue(services.Logger, defaultMaxTasks, defaultTaskCacheSize) // Use constants from task_queue.go
	if err != nil {
		// services.Logger.Fatal().Err(err).Msg("failed to initialize task queue")
		return baseStrategy{}, fmt.Errorf("failed to initialize task queue: %w", err) // Return error
	}
	services.Logger.Info().Msg("task queue initialized within strategy base")

	return baseStrategy{
		ctx:               ctx,
		cancelFunc:        cancel,
		services:          services,
		miner:             miner,
		gpuPool:           gpuPool,
		taskQueue:         taskQueue, // Assign the newly created queue
		logger:            services.Logger.With().Str("strategy", strategyName).Logger(),
		numWorkers:        numWorkers,
		strategyName:      strategyName,
		txParamCache:      paramCache, // Use LRU cache
		needsTaskEvents:   needsTaskEvents,
		sinkTaskSubmitted: sink, // May be nil if needsTaskEvents is false
	}, nil // Return nil error on success
}

func (bs *baseStrategy) gpuWorker(workerId int, gpu *task.GPU, taskHandler taskHandlerFunc) {
	workerLogger := bs.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Logger()
	workerLogger.Info().Msg("started worker")

	newTaskSignal := bs.taskQueue.GetTaskSignal()

	for {
		select {
		case <-bs.ctx.Done(): // Listen to the strategy's context
			workerLogger.Info().Msg("shutting down worker")
			return
		case <-newTaskSignal:
			// Check shutdown condition again after signal
			if bs.ctx.Err() != nil {
				continue
			}

			if !gpu.IsEnabled() {
				workerLogger.Debug().Msg("gpus disabled, skipping task retrieval")
				continue
			}

			ts := bs.taskQueue.GetTask()
			if ts == nil {
				// No task available, maybe another worker got it. Wait for next signal.
				continue
			}

			workerLogger.Debug().Str("task", ts.TaskId.String()).Int("queue_len", bs.taskQueue.Len()).Msg("starting job")
			gpu.SetStatus("Mining")
			taskHandler(workerId, gpu, ts)
			// Task processing finished (success or fail handled within taskHandlerFunc)
			// Status is updated by taskHandlerFunc or remains Mining until next cycle
			workerLogger.Debug().Str("task", ts.TaskId.String()).Msg("finished job processing")
		}
	}

}

// connectToTaskEvents attempts to subscribe to TaskSubmitted events.
// It should only be called by strategies where needsTaskEvents is true.
func (bs *baseStrategy) connectToTaskEvents() {
	if bs.taskEventSub != nil {
		bs.taskEventSub.Unsubscribe()
	}
	bs.logger.Info().Msg("subscribing to TaskSubmitted events...")
	subCtx, cancel := context.WithTimeout(bs.ctx, 15*time.Second)
	defer cancel()
	blockNo, err := bs.services.OwnerAccount.Client.Client.BlockNumber(subCtx)
	if err != nil {
		bs.logger.Error().Err(err).Msg("failed to get current block number for task event subscription")
		// Consider implementing retry logic here or rely on the event loop to retry
		return
	}

	engineContract, err := engine.NewEngine(bs.services.Config.BaseConfig.EngineAddress, bs.services.OwnerAccount.Client.Client)
	if err != nil {
		bs.logger.Error().Err(err).Msg("failed to create engine contract instance for event watching")
		return
	}

	bs.taskEventSub, err = engineContract.WatchTaskSubmitted(&bind.WatchOpts{
		Start:   &blockNo,
		Context: bs.ctx,
	}, bs.sinkTaskSubmitted, nil, nil, nil)
	if err != nil {
		bs.logger.Error().Err(err).Msg("failed to subscribe to TaskSubmitted events")
	} else {
		bs.logger.Info().Msg("subscribed to TaskSubmitted events")
	}
}

// listenForTaskEvents runs in a goroutine for strategies that need TaskSubmitted events.
func (bs *baseStrategy) listenForTaskEvents(handler func(event *engine.EngineTaskSubmitted)) {

	bs.logger.Info().Msg("started task event listener")

	maxBackoff := 30 * time.Second
	currentBackoff := 1 * time.Second

	for {
		select {
		case <-bs.ctx.Done():
			bs.logger.Info().Msg("stopping task event listener")
			return
		case event := <-bs.sinkTaskSubmitted:
			if event == nil {
				continue // Channel closed or nil event
			}
			currentBackoff = 1 * time.Second // Reset backoff on success
			bs.logger.Info().Str("taskid", task.TaskId(event.Id).String()).Str("txHash", event.Raw.TxHash.Hex()).Uint64("block", event.Raw.BlockNumber).Msg("received TaskSubmitted event")
			handler(event) // Call the strategy-specific handler

		case err := <-bs.taskEventSub.Err():
			if err == nil {
				continue
			}
			bs.logger.Warn().Err(err).Msgf("task event subscription error, retrying connection in %s", currentBackoff)
			time.Sleep(currentBackoff)
			currentBackoff = (currentBackoff * 2) + time.Duration(rand.Intn(500))*time.Millisecond
			if currentBackoff > maxBackoff {
				currentBackoff = maxBackoff
			}
			bs.connectToTaskEvents() // Attempt to reconnect
		}
	}
}

// start initializes common strategy components like workers and optionally the event listener.
func (bs *baseStrategy) start(taskHandler taskHandlerFunc, eventHandler func(event *engine.EngineTaskSubmitted)) error {
	bs.logger.Info().Int("workerspergpu", bs.services.Config.NumWorkersPerGPU).Int("gpus", bs.gpuPool.NumGPUs()).Msgf("starting %d workers for %s strategy", bs.numWorkers, bs.strategyName)

	gpus := bs.gpuPool.GetGPUs() // Get a snapshot of GPUs at worker start time
	for i := 0; i < bs.numWorkers; i++ {
		workerId := i
		// Assign GPU in round-robin fashion based on NumWorkersPerGPU
		gpuIndex := workerId / bs.services.Config.NumWorkersPerGPU
		if gpuIndex >= len(gpus) {
			bs.logger.Error().Int("workerId", workerId).Int("gpuIndex", gpuIndex).Int("numGpus", len(gpus)).Msg("logic error: worker assigned to non-existent GPU index")
			continue // Skip starting this worker
		}
		gpu := gpus[gpuIndex]
		bs.Go(func() {
			bs.gpuWorker(workerId, gpu, taskHandler)
		})
	}

	// Start event listener if needed
	if bs.needsTaskEvents {
		if eventHandler == nil {
			return fmt.Errorf("strategy %s needs task events but no handler provided", bs.strategyName)
		}
		bs.connectToTaskEvents()
		bs.Go(func() { bs.listenForTaskEvents(eventHandler) })
	}

	return nil
}

// decodeTransaction retrieves transaction details and decodes parameters using the miner's decoder.
// It utilizes the txParamCache for efficiency.
func (bs *baseStrategy) decodeTransaction(txHash common.Hash) (*SubmitTaskParams, error) {
	// 1. Check cache
	if cachedParams, found := bs.txParamCache.Get(txHash.String()); found {
		if params, ok := cachedParams.(*SubmitTaskParams); ok {
			bs.logger.Debug().Str("txHash", txHash.String()).Msg("using cached task parameters")
			return params, nil
		}
		bs.logger.Warn().Str("txHash", txHash.String()).Msg("invalid type found in txParamCache")
		// Continue to fetch and decode if type assertion failed
	}

	// 2. Fetch transaction if not cached
	txFetchStart := time.Now()
	tx, _, err := bs.services.OwnerAccount.Client.Client.TransactionByHash(bs.ctx, txHash)
	txFetchElapsed := time.Since(txFetchStart)
	if err != nil {
		bs.logger.Error().Err(err).Str("txHash", txHash.String()).Msg("could not get transaction from hash")
		return nil, err
	}
	if tx == nil {
		bs.logger.Error().Str("txHash", txHash.String()).Msg("transaction not found for hash")
		return nil, fmt.Errorf("transaction %s not found", txHash.String())
	}

	bs.logger.Debug().Str("elapsed", txFetchElapsed.String()).Str("hash", txHash.String()).Msg("fetched transaction details")

	// 3. Decode transaction
	params, err := bs.miner.DecodeTaskTransaction(tx)
	if err != nil {
		bs.logger.Error().Err(err).Str("txHash", txHash.String()).Msg("could not decode task transaction")
		return nil, err
	}

	// 4. Store in cache
	bs.txParamCache.Add(txHash.String(), params) // Use Add method of LRU cache
	bs.logger.Debug().Str("txHash", txHash.String()).Msg("cached decoded task parameters")

	return params, nil
}

func (bs *baseStrategy) Stop() {
	bs.stopOnce.Do(func() {
		bs.logger.Info().Msgf("stopping %s strategy workers", bs.strategyName)
		if bs.cancelFunc != nil {
			bs.cancelFunc()
		}
		// Unsubscribe from events if subscribed
		if bs.taskEventSub != nil {
			bs.taskEventSub.Unsubscribe()
			bs.logger.Info().Msg("unsubscribed from task events")
		}
		bs.wg.Wait()
		bs.logger.Info().Msgf("all %s strategy workers stopped", bs.strategyName)
	})
}

// --- BulkMine Strategy ---
// This strategy mines tasks from the storage queue.
// It does NOT listen for TaskSubmitted events.
// It assumes tasks in storage could be from single or bulk submits.

type BulkMineStrategy struct {
	*baseStrategy // Embed pointer
}

func NewBulkMineStrategy(appContext context.Context, services *Services, miner *Miner, gpuPool *GPUPool) (*BulkMineStrategy, error) {
	base, err := newBaseStrategy(appContext, services, miner, gpuPool, "bulkmine", false) // false = does not need task events
	if err != nil {
		return nil, err // Propagate error
	}
	return &BulkMineStrategy{
		baseStrategy: &base, // Store pointer to base
	}, nil
}

func (s *BulkMineStrategy) Name() string { return s.strategyName }

func (s *BulkMineStrategy) Start() error {
	s.logger.Info().Msgf("starting %s strategy", s.Name())
	// Start workers
	err := s.baseStrategy.start(s.handleTask, nil) // No event handler needed
	if err != nil {
		return err
	}
	// Start pulling tasks from storage
	s.Go(s.pullTasksFromStorage)
	return nil
}

// handleTask is the specific task processing logic for BulkMine.
func (s *BulkMineStrategy) handleTask(workerId int, gpu *task.GPU, ts *TaskSubmitted) {
	// assert that the task is not nil
	if ts == nil {
		s.logger.Error().Msg("task is nil")
		return
	}

	workerLogger := s.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", ts.TaskId.String()).Logger()

	// Decode the transaction (uses cache internally)
	params, err := s.decodeTransaction(ts.TxHash)
	if err != nil {
		// Error already logged by decodeTransaction
		s.taskQueue.TaskFailed(ts.TaskId)
		return
	}

	// Solve the task using decoded params
	solveStart := time.Now()
	_, err = s.miner.SolveTask(s.ctx, ts.TaskId, params, gpu, false) // false = not validateOnly
	solveElapsed := time.Since(solveStart)

	if err != nil {
		workerLogger.Error().Err(err).Msg("solve task failed")
		s.taskQueue.TaskFailed(ts.TaskId)
		gpu.SetStatus("Error")
	} else {
		workerLogger.Info().Str("elapsed", solveElapsed.String()).Msg("task solved successfully")
		s.taskQueue.TaskCompleted(ts.TaskId)
		s.gpuPool.AddSolveTime(solveElapsed) // Add successful solve time to pool average
		gpu.SetStatus("Idle")                // Set back to Idle after successful solve
	}
}

// pullTasksFromStorage continuously pulls tasks from the DB and adds them to the queue.
func (s *BulkMineStrategy) pullTasksFromStorage() {
	s.logger.Info().Msg("started task storage puller")

	// TODO: make this configurable or change approach (can we get signal from db task storage)
	// or dont pull tasks if gpus are busy etc.
	ticker := time.NewTicker(20 * time.Millisecond)

	defer ticker.Stop()

	lastEmptyWarning := time.Now()

	poppedCount := 0
	for {
		select {
		case <-s.ctx.Done(): // Use strategy's context
			s.logger.Info().Int("total_popped", poppedCount).Msg("stopping task storage puller")
			return
		case <-ticker.C:
			// Check context again after ticker fires
			if s.ctx.Err() != nil {
				continue
			}

			// Limit queue size before pulling more?
			// if s.taskQueue.Len() > s.taskQueue.maxTasks * 2 { // Example threshold
			// 	s.logger.Debug().Int("queue_size", s.taskQueue.Len()).Msg("task queue is large, pausing storage pull")
			//  time.Sleep(5 * time.Second)
			// 	continue
			// }

			taskId, txHash, err := s.services.TaskStorage.PopTask()
			if err != nil {
				if err == sql.ErrNoRows {
					// Queue is empty, normal condition, wait for next tick
					// we should warn the user there are no tasks to do but only ever n seconds to not flood the logs etc
					if time.Since(lastEmptyWarning) > 10*time.Second {
						s.logger.Warn().Msg("task storage is empty, nothing to do. will sleep for 2 seconds")
						lastEmptyWarning = time.Now()
					}
					time.Sleep(2 * time.Second)
				} else {
					// Log actual DB errors
					s.logger.Error().Err(err).Msg("could not pop task from storage, will retry")
					time.Sleep(5 * time.Second) // Backoff on error
				}
				continue
			}

			// We got a task from storage
			poppedCount++
			s.logger.Info().Str("task", taskId.String()).Int("popped_count", poppedCount).Msg("popped task from storage")

			ts := &TaskSubmitted{
				TaskId: taskId,
				TxHash: txHash,
			}

			// Add to the queue, AddTask handles deduplication and signalling
			added := s.taskQueue.AddTask(ts)
			if !added {
				s.logger.Warn().Str("task", taskId.String()).Msg("popped task from storage was already known/inflight")
			}
			// Optional small delay to prevent overwhelming the queue
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Stop ensures base strategy stop is called.
func (s *BulkMineStrategy) Stop() {
	s.baseStrategy.Stop()
}

// --- AutoMine Strategy ---
// Inherits directly from BulkMineStrategy. Its behavior is identical in terms of
// task processing (pulls from storage, decodes Tx, solves).
// The selection of this strategy signals that BatchTransactionManager should be
// configured to automatically create tasks.

type AutoMineStrategy struct {
	*BulkMineStrategy // Embed BulkMineStrategy as a pointer
}

func NewAutoMineStrategy(appContext context.Context, services *Services, miner *Miner, gpuPool *GPUPool) (*AutoMineStrategy, error) {
	// Initialize the embedded BulkMineStrategy part, but name it "automine"
	bulkStrategy, err := NewBulkMineStrategy(appContext, services, miner, gpuPool)
	if err != nil {
		return nil, err // Propagate error
	}
	bulkStrategy.strategyName = "automine"
	return &AutoMineStrategy{
		BulkMineStrategy: bulkStrategy, // Store the pointer
	}, nil
}

// Name accesses the embedded strategy's method via the pointer.
func (s *AutoMineStrategy) Name() string { return s.strategyName }

// Start accesses the embedded strategy's method via the pointer.
func (s *AutoMineStrategy) Start() error {
	return s.BulkMineStrategy.Start()
}

// Stop accesses the embedded strategy's method via the pointer.
func (s *AutoMineStrategy) Stop() {
	s.BulkMineStrategy.Stop()
}

// HandleTaskSubmitted accesses the embedded strategy's method via the pointer.
func (s *AutoMineStrategy) HandleTaskSubmitted(event *engine.EngineTaskSubmitted) {
	s.BulkMineStrategy.HandleTaskSubmitted(event)
}

// --- SolutionSampler Strategy ---
// Samples tasks from the solution submitted event and adds them to the queue for validation

type SolutionSamplerStrategy struct {
	*baseStrategy                      // Embed pointer
	sampleTicker          *time.Ticker // Ticker for periodic sampling batch processing
	maxTaskSampleSize     int
	tasksSamples          []*TaskSubmitted
	sampleIndex           int
	solutionEventSub      event.Subscription // Hold the subscription to manage errors
	sinkSolutionSubmitted chan *engine.EngineSolutionSubmitted
}

func NewSolutionSamplerStrategy(appCtx context.Context, services *Services, miner *Miner, gpuPool *GPUPool) (*SolutionSamplerStrategy, error) {
	numWorkers := gpuPool.NumGPUs() * services.Config.NumWorkersPerGPU

	sinkSolutionSubmitted := make(chan *engine.EngineSolutionSubmitted, 1024)

	base, err := newBaseStrategy(appCtx, services, miner, gpuPool, "solutionsampler", true) // true = needs task events
	if err != nil {
		return nil, err // Propagate error
	}

	return &SolutionSamplerStrategy{
		// Needs task events to populate the TaskID -> TxHash cache
		baseStrategy:          &base,      // Store pointer to base
		maxTaskSampleSize:     numWorkers, // Sample size based on worker capacity
		tasksSamples:          make([]*TaskSubmitted, 0, numWorkers),
		sinkSolutionSubmitted: sinkSolutionSubmitted,
		solutionEventSub:      nil,
	}, nil
}

// HandleTaskSubmitted caches the TaskID -> TxHash mapping for later lookup when a solution is seen.
func (s *SolutionSamplerStrategy) HandleTaskSubmitted(event *engine.EngineTaskSubmitted) {
	s.taskQueue.CacheTxHash(event.Id, event.Raw.TxHash)
}

func (s *SolutionSamplerStrategy) connectToSolutionSubmittedEvents() {
	if s.solutionEventSub != nil {
		s.solutionEventSub.Unsubscribe()
	}

	s.logger.Info().Msg("subscribing to SolutionSubmitted events...")
	subCtx, cancel := context.WithTimeout(s.ctx, 15*time.Second)
	defer cancel()
	blockNo, err := s.services.OwnerAccount.Client.Client.BlockNumber(subCtx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get current block number for solution event subscription")
		return
	}

	engineContract, err := engine.NewEngine(s.services.Config.BaseConfig.EngineAddress, s.services.OwnerAccount.Client.Client)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create engine contract instance for solution event watching")
		return
	}

	s.solutionEventSub, err = engineContract.WatchSolutionSubmitted(&bind.WatchOpts{
		Start:   &blockNo,
		Context: s.ctx,
	}, s.sinkSolutionSubmitted, nil, nil)

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to subscribe to SolutionSubmitted events")
	} else {
		s.logger.Info().Msg("subscribed to SolutionSubmitted events")
	}

}

func (s *SolutionSamplerStrategy) Name() string { return s.strategyName }

func (s *SolutionSamplerStrategy) Start() error {
	s.logger.Info().Msg("starting SolutionSampler strategy")

	// Start common components (workers, task event listener)
	err := s.baseStrategy.start(s.handleSampledTask, s.HandleTaskSubmitted)
	if err != nil {
		return err
	}

	// Start the goroutine to listen for SolutionSubmitted events
	s.connectToSolutionSubmittedEvents() // Initial connection
	s.Go(s.listenForSolutions)

	// Start the ticker for processing samples periodically
	s.sampleTicker = time.NewTicker(1 * time.Minute) // Process samples every minute
	s.Go(s.periodicSampleProcessor)

	return nil
}

// listenForSolutions handles incoming solution events and adds them to the sample pool.
func (s *SolutionSamplerStrategy) listenForSolutions() {
	s.logger.Info().Msg("started solution event listener")

	maxBackoff := 30 * time.Second
	currentBackoff := 1 * time.Second

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info().Msg("stopping solution event listener")
			return
		case err := <-s.solutionEventSub.Err():
			if err == nil {
				continue
			}
			s.logger.Warn().Err(err).Msgf("solution submitted subscription error, attempting reconnect in %s", currentBackoff)
			time.Sleep(currentBackoff)
			currentBackoff = (currentBackoff * 2) + time.Duration(rand.Intn(500))*time.Millisecond
			if currentBackoff > maxBackoff {
				currentBackoff = maxBackoff
			}
			s.connectToSolutionSubmittedEvents() // Call the reconnect func

		case event := <-s.sinkSolutionSubmitted:
			if event == nil {
				continue
			} // Skip nil events

			// Reset backoff on successful event read
			currentBackoff = 1 * time.Second

			// Get the TxHash from cache (essential for this strategy)
			txHash, found := s.taskQueue.GetCachedTxHash(event.Task)
			if !found {
				s.logger.Warn().Str("task", task.TaskId(event.Task).String()).Msg("solution event received but no TxHash found in cache, skipping sample")
				continue
			}

			ts := &TaskSubmitted{
				TaskId: event.Task,
				TxHash: txHash,
			}

			// Add to sample using reservoir sampling
			s.sampleIndex++
			if len(s.tasksSamples) < s.maxTaskSampleSize {
				s.tasksSamples = append(s.tasksSamples, ts)
				s.logger.Debug().Str("task", ts.TaskId.String()).Int("sample_size", len(s.tasksSamples)).Msg("added solution to sample pool")
			} else {
				j := rand.Intn(s.sampleIndex)
				if j < s.maxTaskSampleSize {
					s.tasksSamples[j] = ts
					s.logger.Debug().Str("task", ts.TaskId.String()).Int("replaced_index", j).Msg("replaced task in sample pool")
				}
			}
		}
	}
}

// periodicSampleProcessor triggers the processing of collected samples.
func (s *SolutionSamplerStrategy) periodicSampleProcessor() {
	s.logger.Info().Msg("started periodic sample processor")

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info().Msg("stopping periodic sample processor")
			return
		case <-s.sampleTicker.C:
			if s.ctx.Err() != nil {
				continue
			}
			s.processSamples()
		}
	}
}

// processSamples copies the current samples and adds them to the task queue
// for the workers to pick up and validate.
func (s *SolutionSamplerStrategy) processSamples() {
	if len(s.tasksSamples) == 0 {
		s.logger.Debug().Msg("no samples collected, skipping processing cycle")
		return
	}

	s.logger.Info().Int("samples", len(s.tasksSamples)).Msg("processing collected solution samples")

	// Create a copy of the samples to avoid race conditions while iterating
	samplesToProcess := make([]*TaskSubmitted, len(s.tasksSamples))
	copy(samplesToProcess, s.tasksSamples)

	// Clear the original sample pool for the next collection period
	s.tasksSamples = s.tasksSamples[:0]
	s.sampleIndex = 0

	// Add the copied samples to the main task queue for validation workers
	addedCount := 0
	for _, ts := range samplesToProcess {
		if s.ctx.Err() != nil { // Check for shutdown during adding
			s.logger.Warn().Msg("shutdown signalled during sample queuing")
			break
		}
		// AddTask handles signalling workers
		if s.taskQueue.AddTask(ts) {
			addedCount++
		} else {
			// Task might already be inflight (e.g., if processed previously and failed?)
			s.logger.Debug().Str("task", ts.TaskId.String()).Msg("sampled task already known/inflight when adding to queue")
		}
	}
	s.logger.Info().Int("added", addedCount).Int("total_sampled", len(samplesToProcess)).Msg("added samples to task queue for validation")
}

// handleSampledTask performs the validation logic for a sampled task.
func (s *SolutionSamplerStrategy) handleSampledTask(workerId int, gpu *task.GPU, ts *TaskSubmitted) {
	workerLogger := s.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", ts.TaskId.String()).Logger()
	workerLogger.Info().Msg("validating sampled task")

	// Decode the transaction (uses cache internally)
	params, err := s.decodeTransaction(ts.TxHash)
	if err != nil {
		// Error logged by decodeTransaction
		s.taskQueue.TaskFailed(ts.TaskId)
		gpu.IncrementErrorCount()
		gpu.SetStatus("Error - Decode")
		return
	}

	// Solve the task in validateOnly mode using decoded params
	solveStart := time.Now()
	ourCidBytes, err := s.miner.SolveTask(s.ctx, ts.TaskId, params, gpu, true) // true = validateOnly
	solveElapsed := time.Since(solveStart)

	if err != nil {
		workerLogger.Error().Err(err).Msg("validation: solve task failed")
		s.taskQueue.TaskFailed(ts.TaskId)
		if gpu.IsEnabled() {
			gpu.SetStatus("Error - Validate")
		}
		return
	}
	if ourCidBytes == nil {
		workerLogger.Error().Msg("validation: solve task did not return a CID")
		s.taskQueue.TaskFailed(ts.TaskId)
		if gpu.IsEnabled() {
			gpu.SetStatus("Error - Validate")
		}
		return
	}
	s.taskQueue.TaskCompleted(ts.TaskId) // Mark validation attempt as complete
	if gpu.IsEnabled() {
		gpu.SetStatus("Idle")
	}
	workerLogger.Info().Str("elapsed", solveElapsed.String()).Msg("validation: task solved locally")

	// Fetch on-chain solution for comparison
	engineContract, err := engine.NewEngine(s.services.Config.BaseConfig.EngineAddress, s.services.OwnerAccount.Client.Client)
	if err != nil {
		workerLogger.Error().Err(err).Msg("validation: failed to create engine contract instance")
		return
	}

	// Use CallOpts with the strategy's context
	callOpts := &bind.CallOpts{Context: s.ctx}
	res, err := engineContract.Solutions(callOpts, ts.TaskId)
	if err != nil {
		workerLogger.Error().Err(err).Msg("validation: error getting on-chain solution info")
		return // Cannot compare
	}

	if res.Blocktime == 0 {
		workerLogger.Warn().Msg("validation: no solution found on-chain (or call failed silently?), cannot compare")
		return
	}

	solversCid := common.Bytes2Hex(res.Cid[:])
	ourCid := common.Bytes2Hex(ourCidBytes)

	workerLogger.Info().Str("our_cid", ourCid).Str("solver_cid", solversCid).Msg("comparing CIDs")

	if ourCid != solversCid {
		workerLogger.Warn().Msg("==================== CID MISMATCH DETECTED =====================")
		workerLogger.Warn().Msgf("  Task ID  : %s", ts.TaskId.String())
		workerLogger.Warn().Msgf("  Our CID  : %s", ourCid)
		workerLogger.Warn().Msgf("  Their CID: %s", solversCid)
		workerLogger.Warn().Msgf("  Solver   : %s", res.Validator.String())
		workerLogger.Warn().Msgf("  Block    : %d", res.Blocktime) // Assuming Blocktime is block number
		workerLogger.Warn().Msg("================================================================")
		// TODO: Add alerting or further action here?
	} else {
		workerLogger.Info().Msg("validation: CID matches on-chain solution")
	}
}

func (s *SolutionSamplerStrategy) Stop() {
	s.logger.Info().Msg("stopping SolutionSampler strategy")
	if s.sampleTicker != nil {
		s.sampleTicker.Stop()
	}
	// Unsubscribe from solution events
	if s.solutionEventSub != nil {
		s.solutionEventSub.Unsubscribe()
		s.logger.Info().Msg("unsubscribed from solution events")
	}
	s.baseStrategy.Stop() // Stop base (workers, task listener)
	s.logger.Info().Msg("SolutionSampler strategy stopped")
}

// --- Listen Strategy ---
// This strategy listens for TaskSubmitted events and adds them directly to the queue.

type ListenStrategy struct {
	*baseStrategy // Embed pointer
}

func NewListenStrategy(appContext context.Context, services *Services, miner *Miner, gpuPool *GPUPool) (*ListenStrategy, error) {
	base, err := newBaseStrategy(appContext, services, miner, gpuPool, "listen", true) // true = needs task events
	if err != nil {
		return nil, err // Propagate error
	}

	return &ListenStrategy{
		baseStrategy: &base, // Store pointer to base
	}, nil
}

func (s *ListenStrategy) Name() string { return s.strategyName }

func (s *ListenStrategy) Start() error {
	s.logger.Info().Msgf("starting %s strategy", s.Name())
	// Start workers and task event listener
	return s.baseStrategy.start(s.handleTask, s.HandleTaskSubmitted)
}

// HandleTaskSubmitted adds the received task event to the processing queue
// and ensures it's recorded in the local task storage.
func (s *ListenStrategy) HandleTaskSubmitted(event *engine.EngineTaskSubmitted) {
	taskId := event.Id
	txHash := event.Raw.TxHash

	// Cast to task.TaskId for String() method and for AddTasks
	taskIdTyped := task.TaskId(taskId)
	taskIdStr := taskIdTyped.String()

	ts := &TaskSubmitted{
		TaskId: taskId, // Store original [32]byte in TaskSubmitted struct
		TxHash: txHash,
	}
	// Add to the queue, AddTask handles deduplication and signalling workers
	added := s.taskQueue.AddTask(ts)
	if !added {
		s.logger.Warn().Str("task", taskIdStr).Msg("received task event but task was already known/inflight in memory queue")
		return
	}

	// zero cost as these are external tasks
	// status 1 = queued
	err := s.services.TaskStorage.AddTaskWithStatus(taskIdTyped, txHash, 0, 1)
	if err != nil {
		s.logger.Error().Err(err).Str("task", taskIdStr).Str("txHash", txHash.Hex()).
			Msg("failed to store received task event in database")
	} else {
		s.logger.Info().Str("task", taskIdStr).Str("txHash", txHash.Hex()).
			Msg("stored received task event in database via AddTask")
	}

}

// handleTask processes a task from the queue (which originated from an event).
func (s *ListenStrategy) handleTask(workerId int, gpu *task.GPU, ts *TaskSubmitted) {
	// Logic is very similar to BulkMineStrategy's handleTask
	if ts == nil {
		s.logger.Error().Msg("task is nil")
		return
	}
	// Cast TaskId here for logging
	workerLogger := s.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", task.TaskId(ts.TaskId).String()).Logger()

	// Decode the transaction (uses cache internally)
	params, err := s.decodeTransaction(ts.TxHash)
	if err != nil {
		// Error already logged by decodeTransaction
		s.taskQueue.TaskFailed(ts.TaskId)
		// TODO: Update task status in DB to indicate decode failure?
		return
	}

	// Solve the task using decoded params
	solveStart := time.Now()
	_, err = s.miner.SolveTask(s.ctx, ts.TaskId, params, gpu, false) // false = not validateOnly
	solveElapsed := time.Since(solveStart)

	if err != nil {
		workerLogger.Error().Err(err).Msg("solve task failed")
		s.taskQueue.TaskFailed(ts.TaskId)
		// TODO: Update task status in DB to indicate solve failure?
		gpu.SetStatus("Error")
	} else {
		workerLogger.Info().Str("elapsed", solveElapsed.String()).Msg("task solved successfully")
		s.taskQueue.TaskCompleted(ts.TaskId)
		// TODO: Update task status in DB to indicate solve success?
		s.gpuPool.AddSolveTime(solveElapsed)
		gpu.SetStatus("Idle")
	}
}

// Stop ensures base strategy stop is called.
func (s *ListenStrategy) Stop() {
	s.baseStrategy.Stop()
}
