package consts

import "time"

const (
	MinMediaWatchInterval      = 1200 * time.Millisecond
	MaxMediaWatchInterval      = 15500 * time.Millisecond
	MaxMediaWatchEndedInterval = MaxMediaWatchInterval + MinMediaWatchInterval

	MediaWatchStartThreshold            = 5000 * time.Millisecond
	MinDurationMediaWatchStartIndicator = 60 * time.Second
)
