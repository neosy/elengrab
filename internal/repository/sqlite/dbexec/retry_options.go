package dbexec

import "time"

type RetryOptions struct {
	MaxRetries int
	Delay      time.Duration
}
