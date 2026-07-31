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
	AddJob(job Job) bool

	// CancelJob removes a job with the given ID from the queue.
	// Returns true if the job was found and removed, false otherwise.
	// If the job is removed, it signals the dispatcher to wake up
	// and re-evaluate available jobs.
	CancelJob(jobID string) bool

	// ActiveWorkers returns the current number of worker goroutines in the pool.
	ActiveWorkers() uint32

	// PoolSize returns the maximum number of workers in the pool.
	PoolSize() uint32
}

type baseWorkerPool struct {
	// logger to use for logging messages
	logger *slog.Logger

	name    string
	options WorkerPoolOptions

	workers map[uint64]Worker

	nextWorkerID atomic.Uint64

	quit chan struct{}
	wg   sync.WaitGroup

	mu   sync.RWMutex
	cond *sync.Cond

	running    atomic.Bool
	terminated atomic.Bool

	semaphore    chan struct{}
	taskStream   chan *task
	jobQueue     *jobQueue
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
		wp.wg.Wait()
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
	wp.wg.Wait()

	if wp.logger != nil {
		wp.logger.Debug(
			"Worker pool manager stopped",
			"name", wp.Name(),
		)
	}
}

// AddJob enqueues a new job for execution.
// job: the job to be executed by a worker.
// Thread-safe; can be called from any goroutine.
func (wp *baseWorkerPool) AddJob(job Job) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.jobQueue.push(job) {
		return false
	}

	wp.cond.Signal() // Notify dispatcher of new job

	return true
}

// CancelJob removes a job with the given ID from the queue.
// Returns true if the job was found and removed, false otherwise.
// If the job is removed, it signals the dispatcher to wake up
// and re-evaluate available jobs.
func (wp *baseWorkerPool) CancelJob(jobID string) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.jobQueue.remove(jobID) {
		wp.cond.Signal()
		return true
	}

	if task, exists := wp.runningTasks[jobID]; exists {
		task.Cancel()
		return true
	}

	return false
}
