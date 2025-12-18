package nworkers

import "time"

type WorkerOptions struct {
	// Worker name
	Name string

	// Interval between repeated runs
	Interval *time.Duration

	// StartAt defines the absolute time when the worker should start.
	// If nil, the worker starts in interval.
	StartAt *time.Time

	// OneShotDelay sets the delay before the first one-shot execution.
	// nil = skip, 0 = run immediately, >0 = run after delay.
	OneShotDelay *time.Duration

	// MaxRuns limits the total number of job executions (one-shot + repeated).
	// 0 = unlimited.
	MaxRuns int
}
