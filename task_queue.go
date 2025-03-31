package main

import (
	"container/list"
	task "gobius/common"
	"sync"
	"time"

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

// TaskQueue manages an in-memory queue of tasks and associated caches.
type TaskQueue struct {
	mu            sync.Mutex
	deque         *list.List           // Holds *TaskSubmitted, newest at front
	inflightTasks map[string]time.Time // Tracks tasks currently being processed: TaskID -> StartTime
	taskHashCache *lru.Cache           // Caches TaskID -> TxHash mapping (primarily for SolutionSampler)
	newTaskSignal chan struct{}        // Signals workers that a new task is available
	maxTasks      int                  // Max size of the deque
	logger        zerolog.Logger
}

// NewTaskQueue creates and initializes a new TaskQueue.
func NewTaskQueue(logger zerolog.Logger, maxTasks int, cacheSize int) (*TaskQueue, error) {
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

	return &TaskQueue{
		deque:         list.New(),
		inflightTasks: make(map[string]time.Time),
		taskHashCache: taskCache,
		newTaskSignal: make(chan struct{}, taskSignalBufferSize), // Buffered channel
		maxTasks:      maxTasks,
		logger:        logger.With().Str("component", "TaskQueue").Logger(),
	}, nil
}

// AddTask adds a new task to the front of the queue and signals workers.
// Returns true if the task was added, false if it was already known or inflight.
func (tq *TaskQueue) AddTask(ts *TaskSubmitted) bool {
	tq.mu.Lock()

	taskIdStr := ts.TaskId.String()

	// Add task to the front
	tq.deque.PushFront(ts)

	//Add to hash cache immediately if needed by samplers even before processing
	tq.taskHashCache.Add(ts.TaskId, ts.TxHash)

	// Evict oldest task if queue exceeds max size
	if tq.deque.Len() > tq.maxTasks {
		tq.deque.Remove(tq.deque.Back())
	}
	tq.mu.Unlock()

	tq.logger.Info().Str("task", taskIdStr).Int("queue_size", tq.deque.Len()).Msg("added task to queue")

	// Signal workers non-blockingly
	select {
	case tq.newTaskSignal <- struct{}{}:
	default:
		// Signal buffer is full, workers are likely busy or backed up.
		// This is usually fine, a worker will pick up the task eventually.
	}

	return true
}

// GetTask retrieves the latest task from the queue for processing.
// Marks the task as inflight. Returns nil if the queue is empty.
func (tq *TaskQueue) GetTask() *TaskSubmitted {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if tq.deque.Len() == 0 {
		return nil // No tasks available
	}

	// Get the latest task from the front
	element := tq.deque.Front()
	ts := element.Value.(*TaskSubmitted)
	tq.deque.Remove(element) // Remove from deque

	taskIdStr := ts.TaskId.String()
	tq.inflightTasks[taskIdStr] = time.Now() // Mark as inflight

	tq.logger.Debug().Str("task", taskIdStr).Int("remaining", tq.deque.Len()).Msg("dequeued task for processing")

	return ts
}

// TaskCompleted marks a task as no longer inflight.
func (tq *TaskQueue) TaskCompleted(taskId task.TaskId) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	taskIdStr := taskId.String()
	delete(tq.inflightTasks, taskIdStr)
	// Keep the task in the hash cache for potential future checks (like solution sampling)
	tq.logger.Debug().Str("task", taskIdStr).Msg("marked task as completed")
}

// TaskFailed marks a task as no longer inflight but potentially requeueable.
// Currently just removes from inflight, could add back to queue if needed.
func (tq *TaskQueue) TaskFailed(taskId task.TaskId) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	taskIdStr := taskId.String()
	delete(tq.inflightTasks, taskIdStr)
	// Optionally add back to the queue with logic to prevent infinite retries
	tq.logger.Warn().Str("task", taskIdStr).Msg("marked task as failed")
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

// InflightCount returns the number of tasks currently being processed.
func (tq *TaskQueue) InflightCount() int {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	return len(tq.inflightTasks)
}
