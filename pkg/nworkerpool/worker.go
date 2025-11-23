package nworkerpool

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
)

type Worker interface {
	Start(ctx context.Context, jobStream chan Job, quit <-chan struct{}, onJobDone func())
	Stop()
	Status() WorkerStatus
	Running() bool
	WorkerId() uint
}

type worker struct {
	logger   *slog.Logger
	workerId uint
	status   atomic.Value
	running  atomic.Bool

	stop chan struct{}
	mu   sync.Mutex
}

func newWorker(
	logger *slog.Logger,
	workerId uint,
) *worker {
	worker := &worker{
		logger:   logger,
		workerId: workerId,
	}

	worker.status.Store(WorkerStatusNone)

	return worker
}

func (w *worker) Start(ctx context.Context, jobStream chan Job, quit <-chan struct{}, onJobDone func()) {
	if !w.running.CompareAndSwap(false, true) {
		if w.logger != nil {
			w.logger.Warn("Worker already running", "workerId", w.workerId)
		}
		return
	}

	// Opening a channel on startup
	w.stop = make(chan struct{})

	w.status.Store(WorkerStatusIdle)

	go func() {
		defer func() {
			w.running.Store(false)
			w.status.Store(WorkerStatusStopped)
			if w.logger != nil {
				w.logger.Debug("Worker stopped", "workerId", w.workerId)
			}
		}()

		for {
			select {
			case <-w.stop:
				return
			case <-quit:
				return
			case job, ok := <-jobStream:
				if !ok {
					if w.logger != nil {
						w.logger.Debug("taskStream closed, stopping worker", "workerId", w.workerId)
					}
					return
				}
				func() {
					w.status.Store(WorkerStatusWorking)
					defer func() {
						w.status.Store(WorkerStatusIdle)
						onJobDone()
					}()

					if w.logger != nil {
						w.logger.Debug("Worker: running job", "workerId", w.workerId)
					}

					job.Execute(ctx, w.workerId)
					if w.logger != nil {
						w.logger.Debug("Worker: job done", "workerId", w.workerId)
					}
				}()
			}
		}
	}()

	if w.logger != nil {
		w.logger.Debug("Worker started", "workerId", w.workerId)
	}
}

func (w *worker) Stop() {
	if !w.running.Load() {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	select {
	case <-w.stop:
		// already closed
	default:
		if w.stop != nil {
			close(w.stop)
		}
	}
}

func (w *worker) Status() WorkerStatus {
	return w.status.Load().(WorkerStatus)
}

func (w *worker) Running() bool {
	return w.running.Load()
}

func (w *worker) WorkerId() uint {
	return w.workerId
}
