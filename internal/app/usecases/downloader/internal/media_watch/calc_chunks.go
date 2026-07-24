package mediawatch

import "time"

func calcChunkCount(duration time.Duration) uint32 {
	if duration <= 0 {
		return 0
	}

	return uint32((duration + chunkDuration - 1) / chunkDuration)
}

func calcRequiredChunkCount(duration time.Duration) uint32 {
	chunks := calcChunkCount(duration)
	if chunks == 0 {
		return 0
	}

	return (chunks*requiredWatchPercent + 100 - 1) / 100
}
