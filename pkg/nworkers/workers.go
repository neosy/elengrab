package nworkers

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
)

type Workers struct {
	logger *slog.Logger

	items []Worker

	running atomic.Bool
	stop    chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
}

func NewWorkers[T any](logger *slog.Logger, deps T, init func(ws *Workers, deps T)) *Workers {
	ws := &Workers{
		logger: logger,
		stop:   make(chan struct{}),
	}

	if init != nil {
		init(ws, deps)
	}

	return ws
}

func (ws *Workers) Add(worker Worker) error {
	if ws.running.Load() {
		if ws.logger != nil {
			ws.logger.Debug("Workers already running")
		}
		return errors.New("workers already running")
	}

	ws.mu.Lock()
	{
		ws.items = append(ws.items, worker)
	}
	ws.mu.Unlock()

	return nil
}

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

func (ws *Workers) StartWorkers(ctx context.Context) error {
	if !ws.running.CompareAndSwap(false, true) {
		if ws.logger != nil {
			ws.logger.Debug("Workers already running")
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
		ws.logger.Debug("Workers running...", "qty", len(ws.items))
	}

	return nil
}

func (ws *Workers) StopWorkers() {
	if !ws.running.Load() {
		return
	}

	ws.mu.Lock()
	{
		select {
		case <-ws.stop:
			return
		default:
			close(ws.stop)
		}
	}
	ws.mu.Unlock()

	ws.Wait()

	if ws.logger != nil {
		ws.logger.Debug("Workers stopped")
	}
}

func (ws *Workers) Wait() {
	ws.wg.Wait()
}
