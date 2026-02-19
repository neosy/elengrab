package nworkerpool

import "time"

const (
	defaultWorkerPoolSize int = 3
	defaultJobQueueCap    int = 100

	defaultDynamicWorkerPoolSize int           = 5
	defaultIdleTime              time.Duration = 5 * time.Second
)
