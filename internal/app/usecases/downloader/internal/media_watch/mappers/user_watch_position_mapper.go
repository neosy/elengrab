package mappers

import (
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (m *Mappers) MapUserWatchEventToWatchPosition(
	event *ddownload.MediaWatchEvent,
	mediaDuration time.Duration,
) *ddownload.MediaUserWatchPosition {
	var userID uuid.UUID
	if event.UserID != nil {
		userID = *event.UserID
	}

	position := &ddownload.MediaUserWatchPosition{
		DownloadID: event.DownloadID,
		UserID:     userID,
		SessionID:  event.SessionID,
		Position:   event.Position,
	}
	position.Normalize(mediaDuration)

	return position
}
