package mediawatch

import "time"

const (
	ChunkDuration = 500 * time.Millisecond

	requiredWatchPercent = 90

	statsUpdateInterval = 15 * time.Second
)
