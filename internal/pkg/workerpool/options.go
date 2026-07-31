package nworkerpool

import (
	"time"
)

type WorkerPoolOption func(*WorkerPoolOptions)

type WorkerPoolOptions struct {
	// maximum number of workers in the pool
	MaxWorkers uint32

	// idle duration before a worker exits
	IdleTime time.Duration
}

// DefaultWorkerPoolOptions returns default options for a worker pool
func DefaultWorkerPoolOptions() WorkerPoolOptions {
	return WorkerPoolOptions{
		MaxWorkers: defaultWorkerMaxWorkers,
	}
}

// DefaultDynamicWorkerPoolOptions returns default options for a dynamic worker pool
func DefaultDynamicWorkerPoolOptions() WorkerPoolOptions {
	return WorkerPoolOptions{
		MaxWorkers: defaultDynamicWorkerMaxWorkers,
		IdleTime:   defaultIdleTime,
	}
}

// ApplyWorkerPoolOptions applies the provided WorkerPoolOption functions to the given WorkerPoolOptions.
func ApplyWorkerPoolOptions(options *WorkerPoolOptions, opts ...WorkerPoolOption) {
	for _, opt := range opts {
		opt(options)
	}
}

// NewWorkerPoolOptions creates a new WorkerPoolOptions instance with the provided options applied.
func NewWorkerPoolOptions(opts ...WorkerPoolOption) WorkerPoolOptions {
	options := DefaultWorkerPoolOptions()

	ApplyWorkerPoolOptions(&options, opts...)

	return options
}

// NewDynamicWorkerPoolOptions creates a new WorkerPoolOptions instance for a dynamic worker pool with the provided options applied.
func NewDynamicWorkerPoolOptions(opts ...WorkerPoolOption) WorkerPoolOptions {
	options := DefaultDynamicWorkerPoolOptions()

	ApplyWorkerPoolOptions(&options, opts...)

	return options
}

// WithMaxWorkers sets the maximum number of workers for a worker pool.
func WithMaxWorkers(maxWorkers uint32) WorkerPoolOption {
	return func(o *WorkerPoolOptions) {
		o.MaxWorkers = maxWorkers
	}
}

// WithIdleTime sets the idle duration for a worker before it exits.
// Only applicable to dynamic worker pools.
func WithIdleTime(idleTime time.Duration) WorkerPoolOption {
	return func(o *WorkerPoolOptions) {
		o.IdleTime = idleTime
	}
}
