package nworkers

import "time"

type WorkerOptions struct {
	// Worker name
	Name string

	// Interval between repeated runs
	Interval *time.Duration

	// Delay before first run
	// If zero, first run happens immediately
	FirstDelay time.Duration
}
