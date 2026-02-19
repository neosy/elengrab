package nworkerpool

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	uformat "github.com/neosy/elengrab/pkg/utils/format"
)

type Worker interface {
	Start(ctx context.Context, task chan *task, quit <-chan struct{}, onJobDone func(jobID string))
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

func (w *worker) Start(ctx context.Context, taskStream chan *task, quit <-chan struct{}, onJobDone func(jobID string)) {
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
			case task, ok := <-taskStream:
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
						onJobDone(task.job.ID())
					}()

					if w.logger != nil {
						w.logger.Debug("Worker: running job", "workerId", w.workerId, "jobName", task.job.Name())
					}

					done := make(chan struct{})
					go func() {
						defer close(done)

						startTime := time.Now()
						task.job.Execute(task.ctx, w.workerId)
						elapsed := time.Since(startTime)

						if w.logger != nil {
							w.logger.Debug(
								"Worker: job done",
								"workerId", w.workerId,
								"jobName", task.job.Name(),
								"elapsed", uformat.DurationFormat(elapsed),
							)
						}
					}()

					select {
					case <-task.ctx.Done():
					case <-done:
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
