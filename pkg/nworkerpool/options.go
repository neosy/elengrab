package nworkerpool

import "time"

// WorkerPoolOptions configures the worker pool behavior.
type WorkerPoolOptions struct {
	// PoolSize  specifies the number of worker goroutines to run.
	// If zero or negative, defaults to defaultWorkerPoolSize.
	PoolSize int
}

type DynamicWorkerPoolOptions struct {
	// PoolSize
	PoolSize int
	//
	IdleTime time.Duration
}

func DefaultDynamicWorkerPoolOptions() DynamicWorkerPoolOptions {
	return DynamicWorkerPoolOptions{
		PoolSize: defaultDynamicWorkerPoolSize,
		IdleTime: defaultIdleTime,
	}
}
