package nworkerpool

import (
	"log/slog"
	"time"
)

type WorkerPoolOption func(*WorkerPoolOptions)

type WorkerPoolOptions struct {
	// logger to use for logging messages
	logger *slog.Logger

	// maximum number of workers in the pool
	MaxWorkers uint32

	// idle duration before a worker exits
	IdleTime time.Duration
}

// DefaultWorkerPoolOptions returns default options for a worker pool
func DefaultWorkerPoolOptions() WorkerPoolOptions {
	return WorkerPoolOptions{
		logger:     nil,
		MaxWorkers: defaultWorkerMaxWorkers,
	}
}

// DefaultDynamicWorkerPoolOptions returns default options for a dynamic worker pool
func DefaultDynamicWorkerPoolOptions() WorkerPoolOptions {
	return WorkerPoolOptions{
		logger:     nil,
		MaxWorkers: defaultDynamicWorkerMaxWorkers,
		IdleTime:   defaultIdleTime,
	}
}

// WorkerPoolOptionLogger sets the logger for a worker pool.
func WorkerPoolOptionLogger(logger *slog.Logger) WorkerPoolOption {
	return func(o *WorkerPoolOptions) {
		o.logger = logger
	}
}

// WorkerPoolOptionMaxWorkers sets the maximum number of workers for a worker pool.
func WorkerPoolOptionMaxWorkers(maxWorkers uint32) WorkerPoolOption {
	return func(o *WorkerPoolOptions) {
		o.MaxWorkers = maxWorkers
	}
}

// WorkerPoolOptionIdleTime sets the idle duration for a worker before it exits.
// Only applicable to dynamic worker pools.
func WorkerPoolOptionIdleTime(idleTime time.Duration) WorkerPoolOption {
	return func(o *WorkerPoolOptions) {
		o.IdleTime = idleTime
	}
}
