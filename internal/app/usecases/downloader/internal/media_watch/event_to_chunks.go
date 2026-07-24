package mediawatch

import (
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaWatch) eventToChunks(event *ddownload.MediaWatchEvent, mediaDuration time.Duration) []*ddownload.MediaWatchChunk {
	if event.Interval <= 0 || event.Position <= 0 {
		return nil
	}

	totalChunks := calcChunkCount(mediaDuration)
	if totalChunks <= 0 {
		return nil
	}

	start := event.Start()
	end := event.Position

	from := calcChunkCount(start)
	to := calcChunkCount(end - 1)

	if to >= totalChunks {
		to = totalChunks - 1
	}

	if from > to {
		return nil
	}

	userID := uuid.Nil
	if event.UserID != nil {
		userID = *event.UserID
	}

	chunks := make([]*ddownload.MediaWatchChunk, 0, to-from+1)
	for i := from; i <= to; i++ {
		chunks = append(chunks, &ddownload.MediaWatchChunk{
			DownloadID: event.DownloadID,
			UserID:     userID,
			ChunkIndex: uint32(i),
			Qty:        1,
		})
	}

	return chunks
}
