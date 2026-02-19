package nworkers

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
)

// Workers manages a set of Worker instances, allowing them to start, stop, and wait for completion.
type Workers struct {
	// optional logger
	logger *slog.Logger

	// registered workers
	items []Worker

	// indicates if the workers are currently running
	running atomic.Bool
	// channel to signal all workers to stop
	stop chan struct{}
	wg   sync.WaitGroup
	mu   sync.Mutex
}

// NewWorkers creates a new Workers manager with optional initialization function.
// T can be used to pass dependencies to the init function.
func NewWorkers[T any](logger *slog.Logger, deps T, init func(logger *slog.Logger, ws *Workers, deps T)) *Workers {
	ws := &Workers{
		logger: logger,
		stop:   make(chan struct{}),
	}

	if init != nil {
		init(logger, ws, deps)
	}

	return ws
}

// Add registers a Worker. Returns an error if workers are already running.
func (ws *Workers) Add(worker Worker) error {
	if ws.running.Load() {
		if ws.logger != nil {
			ws.logger.Warn("Workers already running")
		}
		return errors.New("workers already running")
	}

	ws.mu.Lock()
	ws.items = append(ws.items, worker)
	ws.mu.Unlock()

	return nil
}

// startWorker runs a single worker in its own goroutine and logs start/stop events.
func (ws *Workers) startWorker(ctx context.Context, worker Worker) {
	ws.wg.Go(func() {
		defer func() {
			if ws.logger != nil {
				ws.logger.Debug("Worker stopped", "name", worker.Name())
			}
		}()
		worker.Run(ctx, ws.stop)
	})
	if ws.logger != nil {
		ws.logger.Debug("Worker started", "name", worker.Name())
	}
}

// StartWorker starts a single worker immediately, regardless of whether other workers are currently running.
// If the workers manager is not in a running state, this method transitions it to running before starting the worker.
// The worker is registered in the managed list for tracking and coordination purposes.
// Returns an error only if a concurrent call modified the running state unexpectedly.
func (ws *Workers) StartWorker(ctx context.Context, worker Worker) error {
	// If workers are not running, attempt to transition to running state atomically
	if !ws.running.Load() {
		if !ws.running.CompareAndSwap(false, true) {
			// Concurrent modification: another goroutine started workers between Load and CAS
			if ws.logger != nil {
				ws.logger.Warn("Workers already running")
			}
			return errors.New("workers already running")
		}
	}

	// Register the worker in the managed list under mutex protection
	ws.mu.Lock()
	ws.items = append(ws.items, worker)
	ws.mu.Unlock()

	// Launch the worker goroutine with proper logging and lifecycle management
	ws.startWorker(ctx, worker)

	return nil
}

// StartWorkers starts all registered workers. Returns an error if already running.
func (ws *Workers) StartWorkers(ctx context.Context) error {
	if !ws.running.CompareAndSwap(false, true) {
		if ws.logger != nil {
			ws.logger.Warn("Workers already running")
		}
		return errors.New("workers already running")
	}

	ws.mu.Lock()
	{
		for _, item := range ws.items {
			ws.startWorker(ctx, item)
		}
	}
	ws.mu.Unlock()

	if ws.logger != nil {
		ws.logger.Info("Workers running...", "count", len(ws.items))
	}

	return nil
}

// StopWorkers signals all workers to stop and waits for them to finish.
func (ws *Workers) StopWorkers() {
	if !ws.running.Load() {
		return
	}

	ws.mu.Lock()
	{
		select {
		case <-ws.stop:
			// already stopped
			return
		default:
			close(ws.stop) // signal all workers to stop
		}
	}
	ws.mu.Unlock()

	ws.Wait()

	if ws.logger != nil {
		ws.logger.Debug("Workers stopped")
	}
}

// Wait blocks until all worker goroutines have finished.
func (ws *Workers) Wait() {
	ws.wg.Wait()
}
