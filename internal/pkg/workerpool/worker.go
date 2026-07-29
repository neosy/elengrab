package nworkerpool

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

type Worker interface {
	Start(
		ctx context.Context,
		taskStream chan *task,
		quit <-chan struct{},
		onJobDone func(jobID string),
		onWorkerStop func(workerID uint64),
	)
	StartWithIdleTimeout(
		ctx context.Context,
		idleTime time.Duration,
		taskStream chan *task,
		quit <-chan struct{},
		onJobDone func(jobID string),
		canStopOnIdleTimeout func(workerID uint64) bool,
		onWorkerStop func(workerID uint64),
	)
	Stop()
	Status() WorkerStatus
	Running() bool
	WorkerId() uint64
}

type worker struct {
	logger   *slog.Logger
	workerID uint64
	status   atomic.Value
	running  atomic.Bool

	stop chan struct{}
	mu   sync.Mutex
}

func newWorker(
	logger *slog.Logger,
	workerID uint64,
) *worker {
	worker := &worker{
		logger:   logger,
		workerID: workerID,
	}

	worker.status.Store(WorkerStatusNone)

	return worker
}

func (w *worker) Start(
	ctx context.Context,
	taskStream chan *task,
	quit <-chan struct{},
	onJobDone func(jobID string),
	onWorkerStop func(workerID uint64),
) {
	if !w.running.CompareAndSwap(false, true) {
		if w.logger != nil {
			w.logger.Warn("Worker already running", "workerID", w.workerID)
		}
		return
	}

	// Opening a channel on startup
	w.stop = make(chan struct{})

	w.status.Store(WorkerStatusIdle)

	go func() {
		defer func() {
			if w.logger != nil {
				w.logger.Debug("Worker stopped", "workerID", w.workerID)
			}
			w.running.Store(false)
			w.status.Store(WorkerStatusStopped)
			onWorkerStop(w.workerID)
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
						w.logger.Debug("taskStream closed, stopping worker", "workerID", w.workerID)
					}
					return
				}

				jobRun := func() chan struct{} {
					if w.logger != nil {
						w.logger.Debug("Worker: running job", "workerID", w.workerID, "jobName", task.job.Name())
					}

					done := make(chan struct{})
					go func() {
						defer close(done)

						startTime := time.Now()
						task.job.Execute(task.ctx, w.workerID)
						elapsed := time.Since(startTime)

						if w.logger != nil {
							w.logger.Debug(
								"Worker: job done",
								"workerID", w.workerID,
								"jobName", task.job.Name(),
								"elapsed", uformat.DurationFormat(elapsed),
							)
						}
					}()

					return done
				}

				w.status.Store(WorkerStatusWorking)

				// Run the job and wait for it to complete.
				<-jobRun()

				onJobDone(task.job.ID())
				w.status.Store(WorkerStatusIdle)
			}
		}
	}()
}

func (w *worker) StartWithIdleTimeout(
	ctx context.Context,
	idleTime time.Duration,
	taskStream chan *task,
	quit <-chan struct{},
	onJobDone func(jobID string),
	canStopOnIdleTimeout func(workerID uint64) bool,
	onWorkerStop func(workerID uint64),
) {
	if !w.running.CompareAndSwap(false, true) {
		if w.logger != nil {
			w.logger.Warn("Worker already running", "workerID", w.workerID)
		}
		return
	}

	// Opening a channel on startup
	w.stop = make(chan struct{})

	w.status.Store(WorkerStatusIdle)

	go func() {
		idleTimer := time.NewTimer(idleTime)

		// Reset the timer safely
		idleTimeSafeReset := func() {
			if !idleTimer.Stop() {
				select {
				case <-idleTimer.C:
				default:
				}
			}
			idleTimer.Reset(idleTime)
		}

		defer func() {
			idleTimer.Stop()
			w.running.Store(false)
			w.status.Store(WorkerStatusStopped)
			onWorkerStop(w.workerID)
		}()

		for {
			select {
			case <-w.stop:
				return
			case <-quit:
				return
			case <-idleTimer.C:
				if canStopOnIdleTimeout(w.workerID) {
					return
				}
				idleTimeSafeReset()
			case task, ok := <-taskStream:
				if !ok {
					if w.logger != nil {
						w.logger.Debug("taskStream closed, stopping worker", "workerID", w.workerID)
					}
					return
				}

				jobRun := func() chan struct{} {
					if w.logger != nil {
						w.logger.Debug("Worker: running job", "workerID", w.workerID, "jobName", task.job.Name())
					}

					done := make(chan struct{})
					go func() {
						defer close(done)

						startTime := time.Now()
						task.job.Execute(task.ctx, w.workerID)
						elapsed := time.Since(startTime)

						if w.logger != nil {
							w.logger.Debug(
								"Worker: job done",
								"workerID", w.workerID,
								"jobName", task.job.Name(),
								"elapsed", uformat.DurationFormat(elapsed),
							)
						}
					}()

					return done
				}

				w.status.Store(WorkerStatusWorking)

				// Run the job and wait for it to complete.
				<-jobRun()

				onJobDone(task.job.ID())
				w.status.Store(WorkerStatusIdle)

				idleTimeSafeReset()
			}
		}
	}()
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

func (w *worker) WorkerId() uint64 {
	return w.workerID
}
