package mediawatch

import "time"

const (
	chunkDuration = 100 * time.Millisecond

	requiredWatchPercent = 90

	statsUpdateInterval = 15 * time.Second
)
