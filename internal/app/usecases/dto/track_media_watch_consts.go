package dto

import "time"

const (
	minMediaWatchInterval      = 1200 * time.Millisecond
	maxMediaWatchInterval      = 15500 * time.Millisecond
	maxMediaWatchEndedInterval = maxMediaWatchInterval + minMediaWatchInterval

	mediaWatchStartThreshold            = 5000 * time.Millisecond
	minDurationMediaWatchStartIndicator = 60 * time.Second
)
