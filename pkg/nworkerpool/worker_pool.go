package nworkerpool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type workerPool struct {
	options WorkerPoolOptions

	workers map[uint64]Worker

	nextWorkerID atomic.Uint64

	quit chan struct{}
	wg   sync.WaitGroup

	mu   sync.Mutex
	cond *sync.Cond

	running atomic.Bool

	semaphore    chan struct{}
	taskStream   chan *task
	jobQueue     *jobQueue
	runningTasks map[string]*task
}

// NewWorkerPool creates a new worker pool.
// logger: used for logging events and errors.
// options: optional configuration; nil uses defaults.
// Returns initialized worker pool ready for Start().
func NewWorkerPool(opts ...WorkerPoolOption) *workerPool {
	options := DefaultWorkerPoolOptions()

	for _, opt := range opts {
		opt(&options)
	}

	if options.MaxWorkers == 0 {
		options.MaxWorkers = defaultWorkerMaxWorkers
	}

	wp := &workerPool{
		options:      options,
		workers:      make(map[uint64]Worker, options.MaxWorkers),
		semaphore:    make(chan struct{}, options.MaxWorkers),
		taskStream:   make(chan *task, 1),
		jobQueue:     newJobQueue(defaultJobQueueCap), // Initial but not final queue size
		runningTasks: make(map[string]*task, options.MaxWorkers),
	}

	wp.nextWorkerID.Store(1)
	wp.cond = sync.NewCond(&wp.mu)

	return wp
}

// Start initializes workers and begins job dispatching.
// ctx: passed to each worker for cancellation and timeouts.
// Returns error if manager is already running.
func (wp *workerPool) Start(ctx context.Context) error {
	if !wp.running.CompareAndSwap(false, true) {
		if wp.options.logger != nil {
			wp.options.logger.Warn("Manager already running")
		}
		return errors.New("manager already running")
	}

	// Initialize quit channel and worker pool stop channels.
	wp.quit = make(chan struct{})

	// Create and start all workers
	wp.mu.Lock()
	for range wp.options.MaxWorkers {
		wp.addWorker(ctx)
	}
	wp.mu.Unlock()

	// Bridge ctx cancellation to quit signal and wake dispatcher
	go func(ctx context.Context) {
		select {
		// If manager is already stopping, exit immediately
		case <-wp.quit:
			wp.Stop()
			return
		// On context cancellation: signal shutdown and wake dispatcher
		case <-ctx.Done():
			wp.Stop()
			return
		}
	}(ctx)

	wp.wg.Go(func() {
		// Ensure running flag is cleared when dispatcher exits
		defer func() {
			wp.running.Store(false)
			if wp.options.logger != nil {
				wp.options.logger.Debug("Worker pool manager stopped")
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

			// Extract next job from queue
			job, _ := wp.jobQueue.Pop()
			task := newTask(ctx, job)
			wp.runningTasks[job.ID()] = task

			// Release lock before sending job to worker
			wp.mu.Unlock()

			if job != nil {
				// Claim a worker slot (non-blocking due to prior check)
				select {
				case <-wp.quit:
					return
				case wp.semaphore <- struct{}{}:
				}

				// Dispatch job to an available worker
				select {
				case <-wp.quit:
					return
				case wp.taskStream <- task:
				}
			}
		}
	})

	if wp.options.logger != nil {
		wp.options.logger.Info("Worker pool manager running...", "count", wp.options.MaxWorkers)
	}

	return nil
}

// Stop gracefully shuts down the manager and all workers.
// Idempotent; safe to call multiple times.
// Blocks until all goroutines exit.
func (wp *workerPool) Stop() {
	if !wp.running.CompareAndSwap(true, false) {
		wp.wg.Wait()
		return
	}

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

	if wp.options.logger != nil {
		wp.options.logger.Debug("Worker pool manager stopped")
	}
}

// AddJob enqueues a new job for execution.
// job: the job to be executed by a worker.
// Thread-safe; can be called from any goroutine.
func (wp *workerPool) AddJob(job Job) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.jobQueue.Push(job) {
		return false
	}

	wp.cond.Signal() // Notify dispatcher of new job

	return true
}

// CancelJob removes a job with the given ID from the queue.
// Returns true if the job was found and removed, false otherwise.
// If the job is removed, it signals the dispatcher to wake up
// and re-evaluate available jobs.
func (wp *workerPool) CancelJob(jobID string) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.jobQueue.Remove(jobID) {
		wp.cond.Signal()
		return true
	}

	if task, exists := wp.runningTasks[jobID]; exists {
		task.Cancel()
		return true
	}

	return false
}

// addWorker creates and starts a new worker.
// ctx: context passed to worker for lifecycle control.
// Must be called with wp.mu held.
func (wp *workerPool) addWorker(ctx context.Context) {
	id := wp.nextWorkerID.Load()
	wp.nextWorkerID.Add(1)

	worker := newWorker(wp.options.logger, id)
	wp.workers[id] = worker

	wp.wg.Add(1)
	worker.Start(
		ctx,
		wp.taskStream,
		wp.quit,
		wp.notifyJobDone,
		wp.notifyStopWorker,
	)
}

// notifyStopWorker notifies the manager that a worker is stopping.
func (wp *workerPool) notifyStopWorker(workerID uint64) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.wg.Add(-1)
}

// PoolSize returns the number of worker goroutines in the pool.
func (wp *workerPool) PoolSize() uint32 {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	return wp.options.MaxWorkers
}

// notifyJobDone notifies the manager that a worker finished a job.
func (wp *workerPool) notifyJobDone(jobID string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	delete(wp.runningTasks, jobID)

	select {
	case <-wp.quit:
		return
	case <-wp.semaphore:
	}

	wp.cond.Signal() // Notify dispatcher that a slot is free
}
