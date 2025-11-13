package workerpool

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
	wg   *sync.WaitGroup
}

func newWorker(
	logger *slog.Logger,
	workerId uint,
	wg *sync.WaitGroup,
) *worker {
	worker := &worker{
		logger:   logger,
		workerId: workerId,

		wg:   wg,
		stop: make(chan struct{}),
	}

	worker.status.Store(WorkerStatusNone)

	return worker
}

func (w *worker) Start(ctx context.Context, jobStream chan Job, quit <-chan struct{}, onJobDone func()) {
	if !w.running.CompareAndSwap(false, true) {
		w.logger.Error("Worker already running", "workerId", w.workerId)
		return
	}

	w.status.Store(WorkerStatusIdle)

	w.wg.Go(func() {
		defer func() {
			w.running.Store(false)
			w.status.Store(WorkerStatusStopped)
			w.logger.Debug("Worker stopped", "workerId", w.workerId)
		}()

		for {
			select {
			case <-w.stop:
				return
			case <-quit:
				return
			case job, ok := <-jobStream:
				if !ok {
					w.logger.Debug("taskStream closed, stopping worker", "workerId", w.workerId)
					return
				}
				func() {
					w.status.Store(WorkerStatusWorking)
					defer func() {
						w.status.Store(WorkerStatusIdle)
						onJobDone()
					}()

					w.logger.Debug("Worker: running job", "workerId", w.workerId)
					job.Execute(ctx, w.workerId)
					w.logger.Debug("Worker: job done", "workerId", w.workerId)
				}()
			}
		}
	})

	w.logger.Debug("Worker started", "workerId", w.workerId)
}

func (w *worker) Stop() {
	if !w.running.Load() {
		return
	}

	select {
	case <-w.stop:
		// already closed
	default:
		close(w.stop)
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
