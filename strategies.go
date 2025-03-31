package main

import (
	"context"
	"database/sql"
	"gobius/bindings/engine"
	task "gobius/common"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/rs/zerolog"
)

// MiningStrategy defines the interface for different mining approaches.
type MiningStrategy interface {
	Start() error
	Stop()
	// Name returns the strategy name (e.g., "bulkmine", "solutionsampler")
	Name() string
}

type taskHandlerFunc func(workerId int, gpu *task.GPU, ts *TaskSubmitted)

// --- Base Strategy (Common worker logic) ---
type baseStrategy struct {
	ctx        context.Context // Strategy's operational context
	cancelFunc context.CancelFunc
	services   *Services
	miner      *Miner
	gpuPool    *GPUPool
	taskQueue  *TaskQueue
	logger     zerolog.Logger
	numWorkers int
	wg         sync.WaitGroup
	stopOnce   sync.Once
}

func (b *baseStrategy) Go(f func()) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		f()
	}()
}

func newBaseStrategy(appCtx context.Context, services *Services, miner *Miner, gpuPool *GPUPool, taskQueue *TaskQueue, strategyName string) baseStrategy {
	ctx, cancel := context.WithCancel(appCtx) // Create a derived context for the strategy
	numWorkers := gpuPool.NumGPUs() * services.Config.NumWorkersPerGPU
	return baseStrategy{
		ctx:        ctx,
		cancelFunc: cancel,
		services:   services,
		miner:      miner,
		gpuPool:    gpuPool,
		taskQueue:  taskQueue,
		logger:     services.Logger.With().Str("strategy", strategyName).Logger(),
		numWorkers: numWorkers,
	}
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
				// Optimization: If GPU is known to be disabled, wait for next signal
				// rather than attempting to grab a task immediately.
				// Could add a small sleep here if desired.
				workerLogger.Debug().Msg("GPU disabled, skipping task retrieval")
				continue
			}

			ts := bs.taskQueue.GetTask()
			if ts == nil {
				// No task available, maybe another worker got it. Wait for next signal.
				continue
			}

			workerLogger.Debug().Str("task", ts.TaskId.String()).Int("queue_len", bs.taskQueue.Len()).Msg("starting job")
			gpu.SetStatus("Mining") // Set status before potentially long task

			taskHandler(workerId, gpu, ts) // Execute the strategy-specific handler

			// Task processing finished (success or fail handled within taskHandlerFunc)
			// Status is updated by taskHandlerFunc or remains Mining until next cycle
			workerLogger.Debug().Str("task", ts.TaskId.String()).Msg("finished job processing")
		}
	}

}

// startWorkers launches the common worker goroutines. Specific task handling logic
// is provided by the actual strategy implementation via the taskHandlerFunc.
func (bs *baseStrategy) startWorkers(taskHandler taskHandlerFunc) {
	bs.logger.Info().Int("workerspergpu", bs.services.Config.NumWorkersPerGPU).Int("gpus", bs.gpuPool.NumGPUs()).Msgf("starting %d workers", bs.numWorkers)

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
}

func (bs *baseStrategy) Stop() {
	bs.stopOnce.Do(func() {
		bs.logger.Info().Msg("stopping strategy workers")
		if bs.cancelFunc != nil {
			bs.cancelFunc()
		}
		bs.wg.Wait()
		bs.logger.Info().Msg("all strategy workers stopped")
	})
}

// --- BulkMine Strategy ---

type BulkMineStrategy struct {
	baseStrategy
}

func NewBulkMineStrategy(appContext context.Context, services *Services, miner *Miner, gpuPool *GPUPool, taskQueue *TaskQueue) *BulkMineStrategy {
	return &BulkMineStrategy{
		baseStrategy: newBaseStrategy(appContext, services, miner, gpuPool, taskQueue, "bulkmine"),
	}
}

func (s *BulkMineStrategy) Name() string { return "bulkmine" }

func (s *BulkMineStrategy) Start() error {
	s.logger.Info().Msgf("starting %s strategy", s.Name())

	// Start the common worker pool
	s.baseStrategy.startWorkers(s.handleTask)

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

	// Fetch the full transaction if needed (might already be cached)
	// This potentially blocks, consider timeout or context
	txFetchStart := time.Now()
	tx, _, err := s.services.OwnerAccount.Client.Client.TransactionByHash(s.ctx, ts.TxHash)
	txFetchElapsed := time.Since(txFetchStart)
	if err != nil {
		workerLogger.Error().Err(err).Str("txHash", ts.TxHash.String()).Msg("could not get transaction from hash, failing task")
		s.taskQueue.TaskFailed(ts.TaskId)
		return
	}

	workerLogger.Debug().Str("elapsed", txFetchElapsed.String()).Str("hash", ts.TxHash.String()).Msg("fetched transaction details")

	// Solve the task
	solveStart := time.Now()
	_, err = s.miner.SolveTask(s.ctx, ts.TaskId, tx, gpu, false) // false = not validateOnly
	solveElapsed := time.Since(solveStart)

	if err != nil {
		workerLogger.Error().Err(err).Msg("solve task failed")
		s.taskQueue.TaskFailed(ts.TaskId)
		// Error handling (e.g., incrementing GPU error count) is likely done within SolveTask or GPU client
		// Set status based on error type if possible, otherwise generic Error
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

func (s *BulkMineStrategy) Stop() {
	s.logger.Info().Msg("stopping BulkMine strategy")
	s.baseStrategy.Stop() // Stop the workers first
	s.logger.Info().Msg("BulkMine strategy stopped")
}

// --- SolutionSampler Strategy ---

type SolutionSamplerStrategy struct {
	baseStrategy
	solutionSourceWG      *sync.WaitGroup // WG for the solution event handling goroutine
	sampleTicker          *time.Ticker    // Ticker for periodic sampling batch processing
	maxTaskSampleSize     int
	tasksSamples          []*TaskSubmitted
	sampleIndex           int
	solutionEventSub      event.Subscription // Hold the subscription to manage errors
	sinkSolutionSubmitted chan *engine.EngineSolutionSubmitted
}

func NewSolutionSamplerStrategy(appCtx context.Context, services *Services, miner *Miner, gpuPool *GPUPool, taskQueue *TaskQueue) *SolutionSamplerStrategy {
	numWorkers := gpuPool.NumGPUs() * services.Config.NumWorkersPerGPU

	sinkSolutionSubmitted := make(chan *engine.EngineSolutionSubmitted, 1024)

	return &SolutionSamplerStrategy{
		baseStrategy:          newBaseStrategy(appCtx, services, miner, gpuPool, taskQueue, "solutionsampler"),
		solutionSourceWG:      &sync.WaitGroup{},
		maxTaskSampleSize:     numWorkers, // Sample size based on worker capacity
		tasksSamples:          make([]*TaskSubmitted, 0, numWorkers),
		sinkSolutionSubmitted: sinkSolutionSubmitted,
		solutionEventSub:      nil,
	}
}

func (s *SolutionSamplerStrategy) connectToSolutionSubmittedEvents() {

	var solutionEventSub event.Subscription

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

	solutionEventSub, err = s.services.Engine.Engine.WatchSolutionSubmitted(&bind.WatchOpts{
		Start:   &blockNo,
		Context: s.ctx,
	}, s.sinkSolutionSubmitted, nil, nil)

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to subscribe to SolutionSubmitted events")
	} else {
		s.solutionEventSub = solutionEventSub
		s.logger.Info().Msg("subscribed to SolutionSubmitted events")
	}

}

func (s *SolutionSamplerStrategy) Name() string { return "solutionsampler" }

func (s *SolutionSamplerStrategy) Start() error {
	s.logger.Info().Msg("starting SolutionSampler strategy")

	// Start the common worker pool - for this strategy, workers process sampled tasks
	s.baseStrategy.startWorkers(s.handleSampledTask)

	// Start the goroutine to listen for SolutionSubmitted events
	//s.solutionSourceWG.Add(1)
	go s.listenForSolutions()

	// Start the ticker for processing samples periodically
	s.sampleTicker = time.NewTicker(1 * time.Minute) // Process samples every minute
	//s.solutionSourceWG.Add(1)                        // Track the ticker goroutine as well
	go s.periodicSampleProcessor()

	return nil
}

// listenForSolutions handles incoming solution events and adds them to the sample pool.
func (s *SolutionSamplerStrategy) listenForSolutions() {
	defer s.solutionSourceWG.Done()
	s.logger.Info().Msg("started solution event listener")

	//maxBackoff := 30 * time.Second
	currentBackoff := 1 * time.Second

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info().Msg("stopping solution event listener")
			return
		case err := <-s.solutionEventSub.Err():
			s.logger.Warn().Err(err).Msgf("solution submitted subscription error, attempting reconnect in %s", currentBackoff)
			// Need a way to get the new subscription object after reconnecting
			time.Sleep(currentBackoff)
			s.connectToSolutionSubmittedEvents() // Call the reconnect func

			// currentBackoff *= 2
			// if currentBackoff > maxBackoff {
			// 	currentBackoff = maxBackoff
			// }
		case event := <-s.sinkSolutionSubmitted:
			if event == nil {
				continue
			} // Skip nil events

			// Reset backoff on successful event read
			currentBackoff = 1 * time.Second

			// Get the TxHash from cache (essential for this strategy)
			txHash, found := s.taskQueue.GetCachedTxHash(event.Task)
			if !found {
				// Maybe the task wasn't submitted via our node or cache expired.
				// We could try fetching the tx hash here, but it adds latency/complexity.
				//s.logger.Warn().Str("task", event.Task.String()).Msg("received solution for task not found in cache, cannot sample")
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

	// Fetch the transaction
	txFetchStart := time.Now()
	tx, _, err := s.services.OwnerAccount.Client.Client.TransactionByHash(s.ctx, ts.TxHash)
	txFetchElapsed := time.Since(txFetchStart)
	if err != nil {
		workerLogger.Error().Err(err).Str("txHash", ts.TxHash.String()).Msg("validation: could not get transaction from hash")
		s.taskQueue.TaskFailed(ts.TaskId)
		gpu.IncrementErrorCount()
		gpu.SetStatus("Error - Tx Fetch")
		return
	}

	workerLogger.Debug().Str("elapsed", txFetchElapsed.String()).Msg("validation: fetched transaction details")

	// Solve the task in validateOnly mode
	solveStart := time.Now()
	ourCidBytes, err := s.miner.SolveTask(s.ctx, ts.TaskId, tx, gpu, true) // true = validateOnly
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

	engineContract := s.services.Engine.Engine

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
	s.baseStrategy.Stop()     // Stop workers
	s.solutionSourceWG.Wait() // Wait for event listener and processor
	s.logger.Info().Msg("SolutionSampler strategy stopped")
}
