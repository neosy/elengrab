package mappers

import (
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
)

func (m *Mappers) MapMediaUserWatchPositionDomainToEntity(position *ddownload.MediaUserWatchPosition) (*ewatchevent.MediaUserWatchPosition, error) {
	var sessionID string

	if position.SessionID != uuid.Nil {
		sessionID = position.SessionID.String()
	}

	return &ewatchevent.MediaUserWatchPosition{
		DownloadID: position.DownloadID,
		UserID:     position.UserID,
		SessionID:  sessionID,
		PositionMs: int(position.Position.Milliseconds()),
	}, nil
}

func (m *Mappers) MapMediaUserWatchPositionEntityToDomain(position *ewatchevent.MediaUserWatchPosition) (*ddownload.MediaUserWatchPosition, error) {
	var sessionID uuid.UUID
	if position.SessionID != "" {
		sessionID, _ = uuid.Parse(position.SessionID)
	}

	if position.UserID != uuid.Nil && sessionID != uuid.Nil {
		sessionID = uuid.Nil
	}

	return &ddownload.MediaUserWatchPosition{
		DownloadID: position.DownloadID,
		UserID:     position.UserID,
		SessionID:  sessionID,
		Position:   time.Duration(position.PositionMs) * time.Millisecond,
		CreatedAt:  position.CreatedAt,
		UpdatedAt:  position.UpdatedAt,
	}, nil
}
