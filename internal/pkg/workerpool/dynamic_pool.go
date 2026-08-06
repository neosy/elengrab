package workerpool

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
)

type dynamicWorkerPool struct {
	baseWorkerPool

	stoppingWorkers map[uint64]struct{}

	busyWorkers   uint32
	activeWorkers atomic.Uint32
}

// NewDynamicWorkerPool creates a new dynamic worker pool.
//
// A dynamic worker pool spawns workers on demand: a new worker is
// created when a job arrives and there are no idle workers available.
// Workers automatically exit after being idle for a specified duration.
//
// logger: used for logging events and errors.
// options: optional configuration; nil uses defaults. Supports:
//   - PoolSize: maximum number of concurrent workers (defaults to defaultDynamicWorkerPoolSize)
//   - IdleTime: duration a worker stays alive when idle (defaults to defaultIdleTime)
//
// Returns an initialized *dynamicWorkerPool ready to Start().
func NewDynamicWorkerPool(logger *slog.Logger, name string, opts ...WorkerPoolOption) WorkerPool {
	options := NewDynamicWorkerPoolOptions(opts...)

	if options.MaxWorkers == 0 {
		options.MaxWorkers = defaultDynamicWorkerMaxWorkers
	}

	if options.IdleTime == 0 {
		options.IdleTime = defaultIdleTime
	}

	if logger != nil {
		logger = logger.With("poolName", name)
	}

	wp := &dynamicWorkerPool{
		baseWorkerPool: baseWorkerPool{
			logger:       logger,
			name:         name,
			options:      options,
			workers:      make(map[uint64]Worker, options.MaxWorkers),
			semaphore:    make(chan struct{}, options.MaxWorkers),
			taskStream:   make(chan *task, 1),
			jobQueue:     newJobQueue(defaultJobQueueCap),
			runningTasks: make(map[string]*task, options.MaxWorkers),
		},
		stoppingWorkers: make(map[uint64]struct{}, options.MaxWorkers),
	}

	wp.nextWorkerID.Store(1)
	wp.cond = sync.NewCond(&wp.mu)

	return wp
}

// Start initializes workers and begins job dispatching.
// ctx: passed to each worker for cancellation and timeouts.
// Returns error if manager is already running.
func (wp *dynamicWorkerPool) Start(ctx context.Context) error {
	if wp.terminated.Load() {
		return errors.New("worker pool is terminated")
	}

	if !wp.running.CompareAndSwap(false, true) {
		if wp.logger != nil {
			wp.logger.Warn("Manager already running")
		}
		return errors.New("manager already running")
	}

	// Initialize quit channel and worker pool stop channels.
	wp.quit = make(chan struct{})

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

	wp.globalWG.Go(func() {
		// Ensure running flag is cleared when dispatcher exits
		defer func() {
			wp.running.Store(false)
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
			for wp.jobQueue.len() == 0 || len(wp.semaphore) == cap(wp.semaphore) {
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
			job, _ := wp.jobQueue.pop()
			if job == nil {
				wp.logger.Error(
					"Job queue returned nil job despite being non-empty",
				)
				return
			}

			task := newTask(ctx, job)
			wp.runningTasks[job.ID()] = task

			// Release lock before sending job to worker
			wp.mu.Unlock()

			if task != nil {
				// Claim a worker slot (non-blocking due to prior check)
				select {
				case <-wp.quit:
					return
				case wp.semaphore <- struct{}{}:
					wp.mu.Lock()
					if wp.freeWorkers() == 0 {
						wp.addWorker(ctx)
					}
					wp.busyWorkers++
					wp.mu.Unlock()
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

	if wp.logger != nil {
		wp.logger.Info(
			"Worker pool manager running...",
			"type", "dynamic",
			"activeWorkers", wp.ActiveWorkers(),
			"maxWorkers", wp.options.MaxWorkers,
		)
	}

	return nil
}

// addWorker creates and starts a new worker.
// ctx: context passed to worker for lifecycle control.
// Must be called with wp.mu held.
func (wp *dynamicWorkerPool) addWorker(ctx context.Context) {
	if wp.activeWorkers.Load() >= wp.options.MaxWorkers {
		if wp.logger != nil {
			wp.logger.Warn(
				"Maximum number of workers reached",
				"activeWorkers", wp.ActiveWorkers(),
				"maxWorkers", wp.options.MaxWorkers,
			)
		}

		return
	}

	id := wp.allocateWorkerID()

	worker := newWorker(wp.logger, id)
	wp.workers[id] = worker

	wp.globalWG.Add(1)

	worker.StartWithIdleTimeout(
		ctx,
		wp.options.IdleTime,
		wp.taskStream,
		wp.quit,
		wp.notifyJobDone,
		wp.tryMarkWorkerStopping,
		wp.removeWorker,
	)

	if wp.logger != nil {
		wp.logger.Debug(
			"Worker added to pool",
			"workerID", id,
			"workers", len(wp.workers),
			"activeWorkers", wp.ActiveWorkers(),
			"maxWorkers", wp.options.MaxWorkers,
		)
	}
}

// deleteWorker removes a stopped worker from the pool.
func (wp *dynamicWorkerPool) removeWorker(workerID uint64) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	worker, exists := wp.workers[workerID]
	if !exists {
		if wp.logger != nil {
			wp.logger.Warn(
				"Attempted to remove non-existent worker",
				"workerID", workerID,
			)
		}
		return
	}

	if worker.Status() != WorkerStatusStopped {
		if wp.logger != nil {
			wp.logger.Warn(
				"Attempted to remove non-stopped worker",
				"workerID", workerID,
				"status", worker.Status(),
			)
		}
		return
	}

	delete(wp.workers, workerID)
	delete(wp.stoppingWorkers, workerID)

	wp.releaseWorkerID(workerID)

	wp.globalWG.Done()

	if wp.logger != nil {
		wp.logger.Debug(
			"Worker removed",
			"workerID", workerID,
			"activeWorkers", wp.ActiveWorkers(),
			"maxWorkers", wp.options.MaxWorkers,
		)
	}
}

// ActiveWorkers returns the current number of worker goroutines in the pool.
func (wp *dynamicWorkerPool) ActiveWorkers() uint32 {
	return wp.activeWorkers.Load()
}

// PoolSize returns the maximum number of workers allowed in the dynamic pool.
func (wp *dynamicWorkerPool) PoolSize() uint32 {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.options.MaxWorkers
}

// notifyJobDone notifies the manager that a worker finished a job.
func (wp *dynamicWorkerPool) notifyJobDone(jobID string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.jobsWG.Done()

	select {
	case <-wp.quit:
		return
	case <-wp.semaphore:
		delete(wp.runningTasks, jobID)
		wp.busyWorkers--
		wp.cond.Signal() // Notify dispatcher that a slot is free
		return
	}
}

// tryMarkWorkerStopping checks if a worker can be stopped due to being idle for too long.
func (wp *dynamicWorkerPool) tryMarkWorkerStopping(workerID uint64) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.freeWorkers() > 0 {
		wp.stoppingWorkers[workerID] = struct{}{}
		return true
	}

	return false
}

// allocateWorkerID returns an available worker ID.
// Reuses IDs released by stopped workers before allocating a new ID.
// Must be called with wp.mu held.
func (wp *dynamicWorkerPool) allocateWorkerID() uint64 {
	wp.activeWorkers.Add(1)
	return wp.baseWorkerPool.allocateWorkerID()
}

// releaseWorkerID releases a worker ID for reuse by future workers.
// Must be called with wp.mu held.
func (wp *dynamicWorkerPool) releaseWorkerID(id uint64) {
	wp.activeWorkers.Add(^uint32(0))
	wp.baseWorkerPool.releaseWorkerID(id)
}

// freeWorkers returns the number of workers available to accept new jobs.
// Workers that are busy or scheduled for stopping are excluded.
// Must be called with wp.mu held.
func (wp *dynamicWorkerPool) freeWorkers() uint32 {
	return wp.activeWorkers.Load() - wp.busyWorkers - uint32(len(wp.stoppingWorkers))
}
