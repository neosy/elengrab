package mappers

import (
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
)

func (m *Mappers) MapMediaWatchPositionDomainToEntity(position *ddownload.MediaWatchPosition) (*ewatchevent.MediaWatchPosition, error) {
	var sessionID string

	if position.SessionID != nil {
		sessionID = position.SessionID.String()
	}

	return &ewatchevent.MediaWatchPosition{
		DownloadID: position.DownloadID,
		UserID:     position.UserID,
		SessionID:  sessionID,
		PositionMs: int(position.Position.Milliseconds()),
	}, nil
}

func (m *Mappers) MapMediaWatchPositionEntityToDomain(position *ewatchevent.MediaWatchPosition) (*ddownload.MediaWatchPosition, error) {
	var sessionID *uuid.UUID
	if position.SessionID != "" {
		id, _ := uuid.Parse(position.SessionID)
		if id != uuid.Nil {
			sessionID = &id
		}
	}

	if position.UserID != uuid.Nil && sessionID != nil {
		sessionID = nil
	}

	return &ddownload.MediaWatchPosition{
		DownloadID: position.DownloadID,
		UserID:     position.UserID,
		SessionID:  sessionID,
		Position:   time.Duration(position.PositionMs) * time.Millisecond,
		CreatedAt:  position.CreatedAt,
		UpdatedAt:  position.UpdatedAt,
	}, nil
}
