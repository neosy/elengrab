package workerpool

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
)

type WorkerPool interface {
	// Start initializes workers and begins job dispatching.
	// ctx: passed to each worker for cancellation and timeouts.
	// Returns error if manager is already running.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the manager and all workers.
	// Idempotent; safe to call multiple times.
	// Blocks until all goroutines exit.
	Stop()

	// AddJob enqueues a new job for execution.
	// job: the job to be executed by a worker.
	// Thread-safe; can be called from any goroutine.
	AddJob(job Job) error

	// CancelJob removes a job with the given ID from the queue.
	// Returns true if the job was found and removed, false otherwise.
	// If the job is removed, it signals the dispatcher to wake up
	// and re-evaluate available jobs.
	CancelJob(jobID string) bool

	// ActiveWorkers returns the current number of worker goroutines in the pool.
	ActiveWorkers() uint32

	// PoolSize returns the maximum number of workers in the pool.
	PoolSize() uint32

	// Wait blocks until all pool goroutines have stopped.
	Wait()

	// WaitJobs blocks until all jobs are completed.
	WaitJobs()
}

type baseWorkerPool struct {
	// logger to use for logging messages
	logger *slog.Logger

	name    string
	options WorkerPoolOptions

	workers map[uint64]Worker

	nextWorkerID  atomic.Uint64
	freeWorkerIDs []uint64

	quit chan struct{}
	// all internal goroutines
	globalWG sync.WaitGroup
	jobsWG   sync.WaitGroup

	mu   sync.RWMutex
	cond *sync.Cond

	running    atomic.Bool
	terminated atomic.Bool

	semaphore chan struct{}

	jobQueue *jobQueue

	taskStream   chan *task
	runningTasks map[string]*task
}

// Name returns the name of the worker pool.
func (wp *baseWorkerPool) Name() string {
	return wp.name
}

// Stop gracefully shuts down the manager and all workers.
// Idempotent; safe to call multiple times.
// Blocks until all goroutines exit.
func (wp *baseWorkerPool) Stop() {
	if !wp.running.CompareAndSwap(true, false) {
		wp.globalWG.Wait()
		return
	}

	// Set the terminated flag to true, indicating that the worker pool is shutting down.
	wp.terminated.Store(true)

	// Acquire lock to safely cancel job
	wp.mu.Lock()

	// Signal all components to stop
	close(wp.quit)

	// Signal all jobs to cancel
	for _, task := range wp.runningTasks {
		task.Cancel()
	}

	// Wake up dispatcher if waiting
	wp.cond.Broadcast()

	// Release lock
	wp.mu.Unlock()

	// Wait for dispatcher and workers to finish
	wp.globalWG.Wait()

	if wp.logger != nil {
		wp.logger.Debug(
			"Worker pool manager stopped",
		)
	}
}

// AddJob enqueues a new job for execution.
// job: the job to be executed by a worker.
// Thread-safe; can be called from any goroutine.
func (wp *baseWorkerPool) AddJob(job Job) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.jobQueue.push(job) {
		return ErrJobQueueFull
	}

	wp.jobsWG.Add(1)
	wp.cond.Signal() // Notify dispatcher of new job

	return nil
}

// CancelJob removes a job with the given ID from the queue.
// Returns true if the job was found and removed, false otherwise.
// If the job is removed, it signals the dispatcher to wake up
// and re-evaluate available jobs.
func (wp *baseWorkerPool) CancelJob(jobID string) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.jobQueue.remove(jobID) {
		wp.jobsWG.Done()
		wp.cond.Signal()
		return true
	}

	if task, exists := wp.runningTasks[jobID]; exists {
		task.Cancel()
		return true
	}

	return false
}

// Wait blocks until all submitted jobs have completed.
// It waits for both queued and currently executing jobs to finish.
func (wp *baseWorkerPool) Wait() {
	wp.globalWG.Wait()
}

// WaitJobs blocks until all jobs are completed.
func (wp *baseWorkerPool) WaitJobs() {
	wp.jobsWG.Wait()
}

// allocateWorkerID returns an available worker ID.
// Reuses IDs released by stopped workers before allocating a new ID.
// Must be called with wp.mu held.
func (wp *baseWorkerPool) allocateWorkerID() uint64 {
	n := len(wp.freeWorkerIDs)
	if n > 0 {
		id := wp.freeWorkerIDs[n-1]
		wp.freeWorkerIDs = wp.freeWorkerIDs[:n-1]
		return id
	}

	return wp.nextWorkerID.Add(1) - 1
}

// releaseWorkerID releases a worker ID for reuse by future workers.
// Must be called with wp.mu held.
func (wp *baseWorkerPool) releaseWorkerID(id uint64) {
	wp.freeWorkerIDs = append(wp.freeWorkerIDs, id)
}
