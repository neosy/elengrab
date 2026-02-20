package nworkerpool

import "time"

const (
	// default maximum number of workers in a static pool
	defaultWorkerMaxWorkers uint32 = 3

	// default capacity of the job queue
	defaultJobQueueCap uint32 = 100
)

const (
	// default maximum number of workers in a dynamic pool
	defaultDynamicWorkerMaxWorkers uint32 = 5

	// default idle time before a dynamic worker exits
	defaultIdleTime time.Duration = 5 * time.Second
)
