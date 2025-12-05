package nworkerpool

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
)

const (
	defaultWorkerPoolSize int = 3
	defaultJobQueueCap    int = 100
)

// WorkerPool represents a pool of workers that can execute jobs concurrently.
// It provides methods to start and stop the pool and query the number of workers.
type WorkerPool interface {
	// Start begins execution of the worker pool.
	// ctx: context for cancellation and worker lifecycle.
	// Returns error if already running.
	Start(ctx context.Context) error

	// Stop gracefully shuts down all workers and the pool
	// Idempotent; safe to call multiple times.
	Stop()

	// PoolSize returns the number of worker goroutines in the pool.
	PoolSize() int
}

// WorkerPoolOptions configures the worker pool behavior.
type WorkerPoolOptions struct {
	// PoolSize  specifies the number of worker goroutines to run.
	// If zero or negative, defaults to defaultWorkerPoolSize.
	PoolSize int
}

type workerPool struct {
	logger   *slog.Logger
	workers  []Worker
	poolSize int

	quit    chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
	running atomic.Bool
	cond    *sync.Cond

	semaphore chan struct{}
	jobStream chan Job
	jobQueue  jobQueue
}

// NewWorkerPool creates a new worker pool.
// logger: used for logging events and errors.
// options: optional configuration; nil uses defaults.
// Returns initialized worker pool ready for Start().
func NewWorkerPool(logger *slog.Logger, options *WorkerPoolOptions) *workerPool {
	var poolSize = defaultWorkerPoolSize
	if options != nil && options.PoolSize > 0 {
		poolSize = options.PoolSize
	}

	wp := &workerPool{
		logger:    logger,
		workers:   make([]Worker, 0, poolSize),
		poolSize:  poolSize,
		quit:      make(chan struct{}),
		semaphore: make(chan struct{}, poolSize),
		jobStream: make(chan Job, 1),
		jobQueue:  newJobQueue(defaultJobQueueCap), // Initial but not final queue size
	}

	wp.cond = sync.NewCond(&wp.mu)

	return wp
}

// Start initializes workers and begins job dispatching.
// ctx: passed to each worker for cancellation and timeouts.
// Returns error if manager is already running.
func (wp *workerPool) Start(ctx context.Context) error {
	if !wp.running.CompareAndSwap(false, true) {
		if wp.logger != nil {
			wp.logger.Warn("Manager already running")
		}
		return errors.New("manager already running")
	}

	// Create and start all workers
	wp.mu.Lock()
	for range wp.poolSize {
		wp.addWorker(ctx)
	}
	wp.mu.Unlock()

	// Bridge ctx cancellation to quit signal and wake dispatcher
	go func(ctx context.Context) {
		select {
		// If manager is already stopping, exit immediately
		case <-wp.quit:
			return
		// On context cancellation: signal shutdown and wake dispatcher
		case <-ctx.Done():
			wp.mu.Lock()
			{
				select {
				case <-wp.quit:
				default:
					close(wp.quit)
				}
				wp.cond.Broadcast() // Wake dispatcher from cond.Wait()
			}
			wp.mu.Unlock()
			return
		}
	}(ctx)

	wp.wg.Go(func() {
		// Ensure running flag is cleared when dispatcher exits
		defer func() {
			wp.running.Store(false)
			if wp.logger != nil {
				wp.logger.Debug("Worker pool manager stopped")
			}
		}()

		for {
			// Fast-path: exit immediately if quit signal received
			select {
			case <-wp.quit:
				return
			default:
			}

			// Acquire lock to safely inspect queue and semaphore
			wp.mu.Lock()
			{
				// Wait until there is at least one job AND a free worker slot
				for wp.jobQueue.Len() == 0 || len(wp.semaphore) == cap(wp.semaphore) {
					select {
					// Stop loop if quit signal is received during wait
					case <-wp.quit:
						wp.mu.Unlock()
						return
					default:
						// Release the mutex lock and sleep until notified (job added or slot freed)
						wp.cond.Wait()
					}
				}
			}
			// Release lock before sending job to worker
			wp.mu.Unlock()

			// Extract next job from queue
			job, _ := wp.jobQueue.Pop()

			if job != nil {
				// Claim a worker slot (non-blocking due to prior check)
				select {
				case <-wp.quit:
				case wp.semaphore <- struct{}{}:
				}

				// Dispatch job to an available worker
				select {
				case <-wp.quit:
				case wp.jobStream <- job:
				}
			}
		}
	})

	if wp.logger != nil {
		wp.logger.Debug("Worker pool manager running...", "count", wp.poolSize)
	}

	return nil
}

// Stop gracefully shuts down the manager and all workers.
// Idempotent; safe to call multiple times.
// Blocks until all goroutines exit.
func (wp *workerPool) Stop() {
	if !wp.running.CompareAndSwap(true, false) {
		return
	}

	wp.mu.Lock()
	{
		// Signal all components to stop
		close(wp.quit)

		// Wake up dispatcher if waiting
		wp.cond.Broadcast()
	}
	wp.mu.Unlock()

	// Wait for dispatcher and workers to finish
	wp.wg.Wait()
}

// AddJob enqueues a new job for execution.
// job: the job to be executed by a worker.
// Thread-safe; can be called from any goroutine.
func (wp *workerPool) AddJob(job Job) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.jobQueue.Push(job)
	wp.cond.Broadcast() // Notify dispatcher of new job
}

// addWorker creates and starts a new worker.
// ctx: context passed to worker for lifecycle control.
// Must be called with wp.mu held.
func (wp *workerPool) addWorker(ctx context.Context) {
	worker := newWorker(wp.logger, uint(len(wp.workers)))
	wp.workers = append(wp.workers, worker)

	worker.Start(
		ctx,
		wp.jobStream,
		wp.quit,
		wp.notifyJobDone,
	)
}

// PoolSize returns the number of worker goroutines in the pool.
func (wp *workerPool) PoolSize() int {
	return wp.poolSize
}

// notifyJobDone notifies the manager that a worker finished a job.
func (wp *workerPool) notifyJobDone() {
	select {
	case <-wp.quit:
		return
	case <-wp.semaphore:
	}
	wp.mu.Lock()
	wp.cond.Broadcast() // Notify dispatcher that a slot is free
	wp.mu.Unlock()
}
