package nworkers

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// WorkerJob represents a single unit of work to be executed by a Worker
type WorkerJob interface {
	Execute(ctx context.Context) error
}

// Worker defines the interface for a worker that can run jobs periodically.
type Worker interface {
	// Starts the worker loop
	Run(ctx context.Context, stop <-chan struct{}) error
	// Returns worker name
	Name() string
}

// worker is a concrete implementation of Worker.
type worker struct {
	// indicates if the worker is currently running
	running atomic.Bool
	// the job to execute
	job WorkerJob

	// options
	optons WorkerOptions
}

// NewWorker creates a new Worker with the given job and options.
func NewWorker(job WorkerJob, options *WorkerOptions) Worker {
	w := &worker{
		job: job,
	}

	if options != nil {
		w.optons = *options
	}

	return w
}

// Run starts the worker execution loop.
//
// The job may be executed:
//   - once after OneShotDelay, if set;
//   - at StartAt and then repeatedly with Interval, if both are set;
//   - repeatedly with Interval, if Interval is set without StartAt.
//
// Run returns immediately if no execution options are configured.
//
// The worker stops when the context is cancelled, the stop channel is signaled,
// or MaxRuns (if set) is reached.
func (w *worker) Run(ctx context.Context, stop <-chan struct{}) error {
	if !w.running.CompareAndSwap(false, true) {
		return errors.New("worker already running")
	}
	defer w.running.Store(false)

	var (
		timer         *time.Timer
		timerC        <-chan time.Time = nil
		ticker        *time.Ticker
		tickerC       <-chan time.Time = nil
		startAtTimer  *time.Timer
		startAtTimerC <-chan time.Time = nil
	)

	defer func() {
		if timer != nil {
			timer.Stop()
		}
		if ticker != nil {
			ticker.Stop()
		}
		if startAtTimer != nil {
			startAtTimer.Stop()
		}
	}()

	// One-shot execution after an optional initial delay.
	if w.optons.OneShotDelay != nil {
		timer = time.NewTimer(*w.optons.OneShotDelay)
		timerC = timer.C
	}

	// Periodic execution
	if w.optons.StartAt == nil && w.optons.Interval != nil {
		ticker = time.NewTicker(*w.optons.Interval)
		tickerC = ticker.C
	}

	if w.optons.StartAt != nil && w.optons.Interval != nil {
		startAt := *w.optons.StartAt
		now := time.Now().UTC()

		if startAt.Before(now) {
			// move startAt forward to the next valid interval
			elapsed := now.Sub(startAt)
			steps := elapsed / *w.optons.Interval
			startAt = startAt.Add((steps + 1) * (*w.optons.Interval))
		}

		delay := time.Until(startAt)
		startAtTimer = time.NewTimer(delay)
		startAtTimerC = startAtTimer.C
	}

	// Nothing to execute: neither delayed nor periodic execution is configured.
	if timerC == nil && tickerC == nil && startAtTimerC == nil {
		return nil
	}

	runs := 0
	for {
		select {
		case <-ctx.Done():
			// Context cancellation always terminates the worker
			return nil
		case <-stop:
			// Explicit stop signal terminates the worker.
			return nil
		case <-timerC:
			// Initial one-shot execution.
			timerC = nil
			w.job.Execute(ctx)
			runs++

			if tickerC == nil && startAtTimerC == nil {
				return nil
			}

			if w.optons.MaxRuns > 0 && runs >= w.optons.MaxRuns {
				return nil
			}
		case <-startAtTimerC:
			startAtTimerC = nil
			w.job.Execute(ctx)
			runs++

			if w.optons.MaxRuns > 0 && runs >= w.optons.MaxRuns {
				return nil
			}

			if ticker == nil {
				ticker = time.NewTicker(*w.optons.Interval)
				tickerC = ticker.C
			}
		case <-tickerC:
			// Periodic execution.
			w.job.Execute(ctx)
			runs++
			if w.optons.MaxRuns > 0 && runs >= w.optons.MaxRuns {
				return nil
			}
		}
	}
}

// Name returns the worker's name.
func (w *worker) Name() string {
	return w.optons.Name
}
