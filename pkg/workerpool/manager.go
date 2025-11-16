package workerpool

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
)

const (
	workerCountDefault int = 3
)

// Manager defines the interface for controlling the worker pool lifecycle.
type Manager interface {
	// Start begins execution of the worker pool.
	// ctx: context for cancellation and worker lifecycle.
	// Returns error if already running.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the manager and all workers.
	// Idempotent; safe to call multiple times.
	Stop()

	// WorkerCount returns the number of worker goroutines in the pool.
	WorkerCount() int
}

// ManagerOptions configures the worker pool behavior.
type ManagerOptions struct {
	// WorkerCount specifies the number of worker goroutines to run.
	// If zero or negative, defaults to maxWorkersDefault.
	WorkerCount int
}

type manager struct {
	logger      *slog.Logger
	workers     []Worker
	workerCount int

	quit    chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
	running atomic.Bool
	cond    *sync.Cond

	semaphore chan struct{}
	jobStream chan Job
	jobQueue  jobQueue
}

// NewManager creates a new worker pool manager.
// logger: used for logging events and errors.
// options: optional configuration; nil uses defaults.
// Returns initialized manager ready for Start().
func NewManager(logger *slog.Logger, options *ManagerOptions) *manager {
	var workerCount = workerCountDefault
	if options != nil && options.WorkerCount > 0 {
		workerCount = options.WorkerCount
	}

	m := &manager{
		logger:      logger,
		workers:     make([]Worker, 0),
		workerCount: workerCount,
		quit:        make(chan struct{}),
		semaphore:   make(chan struct{}, workerCount),
		jobStream:   make(chan Job, 1),
		jobQueue:    *newJobQueue(100),
	}

	m.cond = sync.NewCond(&m.mu)

	return m
}

// Start initializes workers and begins job dispatching.
// ctx: passed to each worker for cancellation and timeouts.
// Returns error if manager is already running.
func (m *manager) Start(ctx context.Context) error {
	if !m.running.CompareAndSwap(false, true) {
		err := errors.New("manager already running")
		m.logger.Error("Manager already running")
		return err
	}

	// Create and start all workers
	for range m.workerCount {
		m.addWorker(ctx)
	}

	// Bridge ctx cancellation to quit signal and wake dispatcher
	go func(ctx context.Context) {
		select {
		// If manager is already stopping, exit immediately
		case <-m.quit:
			return
		// On context cancellation: signal shutdown and wake dispatcher
		case <-ctx.Done():
			m.mu.Lock()
			select {
			case <-m.quit:
			default:
				close(m.quit)
			}
			m.cond.Broadcast() // Wake dispatcher from cond.Wait()
			m.mu.Unlock()
			return
		}
	}(ctx)

	m.wg.Go(func() {
		// Ensure running flag is cleared when dispatcher exits
		defer func() {
			m.running.Store(false)
			m.logger.Debug("Worker pool manager stopped")
		}()

		for {
			// Fast-path: exit immediately if quit signal received
			select {
			case <-m.quit:
				return
			default:
			}

			// Acquire lock to safely inspect queue and semaphore
			m.mu.Lock()

			// Wait until there is at least one job AND a free worker slot
			for m.jobQueue.Len() == 0 || len(m.semaphore) == cap(m.semaphore) {
				select {
				// Stop loop if quit signal is received during wait
				case <-m.quit:
					m.mu.Unlock()
					return
				default:
					// Release lock and sleep until notified (job added or slot freed)
					m.cond.Wait()
				}
			}

			// Claim a worker slot (non-blocking due to prior check)
			m.semaphore <- struct{}{}

			// Extract next job from queue
			job, _ := m.jobQueue.Pop()

			// Release lock before sending job to worker
			m.mu.Unlock()

			// Dispatch job to an available worker
			m.jobStream <- job
		}
	})

	m.logger.Debug("Worker pool manager running...")

	return nil
}

// Stop gracefully shuts down the manager and all workers.
// Idempotent; safe to call multiple times.
// Blocks until all goroutines exit.
func (m *manager) Stop() {
	if !m.running.CompareAndSwap(true, false) {
		return
	}

	// Signal all components to stop
	close(m.quit)

	// Wake up dispatcher if waiting
	m.mu.Lock()
	m.cond.Broadcast()
	m.mu.Unlock()

	// Wait for dispatcher and workers to finish
	m.wg.Wait()
}

// AddJob enqueues a new job for execution.
// job: the job to be executed by a worker.
// Thread-safe; can be called from any goroutine.
func (m *manager) AddJob(job Job) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobQueue.Push(job)
	m.cond.Broadcast() // Notify dispatcher of new job
}

// addWorker creates and starts a new worker.
// ctx: context passed to worker for lifecycle control.
// Must be called with m.mu held.
func (m *manager) addWorker(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	worker := newWorker(
		m.logger,
		uint(len(m.workers)),
		&m.wg,
	)
	m.workers = append(m.workers, worker)

	worker.Start(
		ctx,
		m.jobStream,
		m.quit,
		func() {
			<-m.semaphore

			m.mu.Lock()
			m.cond.Broadcast() // Notify dispatcher that a slot is free
			m.mu.Unlock()
		},
	)
}

// WorkerCount returns the number of worker goroutines in the pool.
func (m *manager) WorkerCount() int {
	return m.workerCount
}
