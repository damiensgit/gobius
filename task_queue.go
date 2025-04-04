package main

import (
	"container/list"
	task "gobius/common"
	"sync"

	"context"

	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog"
)

const (
	defaultMaxTasks      = 50    // Maximum number of tasks to store in the deque
	defaultTaskCacheSize = 10000 // Default size for the task hash cache
	taskSignalBufferSize = 1     // Buffer size for the new task signal channel
)

// TaskSubmitted holds information about a submitted task relevant for the queue.
type TaskSubmitted struct {
	TaskId task.TaskId
	TxHash common.Hash
	// Could add BlockNumber, Timestamp here if needed later
}

// TaskQueue manages tasks waiting for GPU processing.
type TaskQueue struct {
	logger         zerolog.Logger
	deque          *list.List               // Stores *TaskSubmitted
	maxTasks       int                      // Max size of the deque
	taskMap        map[string]*list.Element // TaskId.String() -> Element
	taskHashCache  *lru.Cache               // Cache TaskId -> TxHash
	mu             sync.Mutex
	newTaskSignal  chan struct{} // Buffered channel (size 1) to signal workers AND block producer
	closed         bool
	enableEviction bool            // Flag to control eviction behavior
	ctx            context.Context // Add context for shutdown signal in AddTask
}

// NewTaskQueue creates and initializes a new TaskQueue.
func NewTaskQueue(ctx context.Context, logger zerolog.Logger, maxTasks int, cacheSize int, enableEviction bool) (*TaskQueue, error) {
	if maxTasks <= 0 {
		maxTasks = defaultMaxTasks
	}
	if cacheSize <= 0 {
		cacheSize = defaultTaskCacheSize
	}

	taskCache, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}

	tq := &TaskQueue{
		logger:         logger.With().Str("component", "TaskQueue").Logger(),
		deque:          list.New(),
		maxTasks:       maxTasks,
		taskMap:        make(map[string]*list.Element),
		taskHashCache:  taskCache,
		newTaskSignal:  make(chan struct{}, 1), // Explicitly buffer size 1
		closed:         false,
		enableEviction: enableEviction,
		ctx:            ctx, // Store context
	}
	return tq, nil
}

// AddTask adds a task to the queue. Checks capacity if eviction is disabled.
// Performs a BLOCKING send to signal workers, providing backpressure.
// Returns true if the task was added and signaled, false otherwise.
func (tq *TaskQueue) AddTask(ts *TaskSubmitted) bool {
	tq.mu.Lock()

	if tq.closed {
		tq.mu.Unlock()
		return false
	}

	taskIdStr := ts.TaskId.String()

	// Check duplicates
	if _, exists := tq.taskMap[taskIdStr]; exists {
		tq.logger.Warn().Str("task", taskIdStr).Msg("task already in queue")
		tq.mu.Unlock()
		return false
	}

	// Check capacity ONLY if eviction is disabled
	if !tq.enableEviction && tq.deque.Len() >= tq.maxTasks {
		tq.logger.Debug().Int("deque_len", tq.deque.Len()).Int("max_tasks", tq.maxTasks).Msg("Queue full and eviction disabled, not adding task")
		tq.mu.Unlock()
		return false // Indicate queue is full
	}

	// Evict oldest if enabled and queue is full
	if tq.enableEviction && tq.deque.Len() >= tq.maxTasks {
		tq.EvictOldest()
	}

	// Add the new task
	element := tq.deque.PushFront(ts)
	tq.taskMap[taskIdStr] = element
	tq.CacheTxHash(ts.TaskId, ts.TxHash)
	qLen := tq.deque.Len()

	tq.mu.Unlock() // Unlock *before* blocking send

	tq.logger.Debug().Str("task", taskIdStr).Int("queue_len", qLen).Msg("added task to deque, attempting blocking signal")

	// Blocking send on size-1 channel. Waits if buffer is full (worker is busy).
	// Also selects on context cancellation to allow shutdown while blocked.
	select {
	case tq.newTaskSignal <- struct{}{}:
		tq.logger.Debug().Str("task", taskIdStr).Msg("signalled worker successfully")
		return true // Signal sent, task effectively added
	case <-tq.ctx.Done():
		tq.logger.Info().Str("task", taskIdStr).Msg("add task interrupted by context cancellation")
		// We added the task to the deque but couldn't signal.
		// Need to potentially remove it to maintain consistency?
		tq.mu.Lock()
		if elem, ok := tq.taskMap[taskIdStr]; ok {
			// Check it wasn't processed in the meantime
			delete(tq.taskMap, taskIdStr)
			tq.deque.Remove(elem)
			//tq.logger.Warn().Str("task", taskIdStr).Msg("removed task from deque due to context cancellation during signal send")
		}
		tq.mu.Unlock()
		return false // Indicate failure due to shutdown
	}
}

// GetTask retrieves the next available task from the queue.
// Returns nil if the queue is empty.
func (tq *TaskQueue) GetTask() *TaskSubmitted {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if tq.deque.Len() == 0 {
		return nil
	}

	element := tq.deque.Back()
	tq.deque.Remove(element)
	ts := element.Value.(*TaskSubmitted)
	taskIdStr := ts.TaskId.String()

	delete(tq.taskMap, taskIdStr)

	tq.logger.Debug().Str("task", taskIdStr).Int("queue_len", tq.deque.Len()).Msg("retrieved task from queue")

	return ts
}

// Stop signals the queue to stop accepting new tasks.
func (tq *TaskQueue) Stop() {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if tq.closed {
		return
	}
	tq.closed = true
	close(tq.newTaskSignal)
	tq.logger.Info().Msg("task queue stopped")
}

// TaskCompleted marks a task as no longer inflight.
func (tq *TaskQueue) TaskCompleted(taskId task.TaskId) {
	tq.logger.Debug().Str("task", taskId.String()).Msg("marked task as completed")
}

// TaskFailed marks a task as no longer inflight but potentially requeueable.
func (tq *TaskQueue) TaskFailed(taskId task.TaskId) {
	tq.logger.Warn().Str("task", taskId.String()).Msg("marked task as failed")
}

// GetTaskSignal returns the channel workers listen on for new task notifications.
func (tq *TaskQueue) GetTaskSignal() <-chan struct{} {
	return tq.newTaskSignal
}

// GetCachedTxHash retrieves the cached transaction hash for a given Task ID.
// Useful for strategies like SolutionSampler. Returns the hash and true if found.
func (tq *TaskQueue) GetCachedTxHash(taskId task.TaskId) (common.Hash, bool) {
	val, found := tq.taskHashCache.Get(taskId)
	if !found {
		return common.Hash{}, false
	}
	hash, ok := val.(common.Hash)
	if !ok {
		// This indicates a programming error (wrong type stored in cache)
		tq.logger.Error().Str("task", taskId.String()).Msg("invalid type found in task hash cache")
		return common.Hash{}, false
	}
	return hash, true
}

// CacheTxHash explicitly adds a task ID and tx hash to the cache.
// This is an alias for AddToCache, potentially used by specific strategies.
func (tq *TaskQueue) CacheTxHash(taskId task.TaskId, txHash common.Hash) {
	tq.AddToCache(taskId, txHash)
}

// AddToCache explicitly adds a task ID and tx hash to the cache.
// Used when a task isn't added via AddTask but needs caching (e.g., SolutionSampler).
func (tq *TaskQueue) AddToCache(taskId task.TaskId, txHash common.Hash) {
	// Add even if it might already exist, LRU handles eviction
	tq.taskHashCache.Add(taskId, txHash)
}

// Len returns the current number of tasks waiting in the queue.
func (tq *TaskQueue) Len() int {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	return tq.deque.Len()
}

// EvictOldest removes the oldest task from the queue.
func (tq *TaskQueue) EvictOldest() {
	// Called under lock
	oldest := tq.deque.Back()
	if oldest != nil {
		ts := oldest.Value.(*TaskSubmitted)
		tq.deque.Remove(oldest)
		taskIdStr := ts.TaskId.String()
		delete(tq.taskMap, taskIdStr)
		tq.taskHashCache.Remove(ts.TaskId)
		// No need to consume signal here
		tq.logger.Debug().Str("task", taskIdStr).Msg("evicted oldest task")
	}
}
