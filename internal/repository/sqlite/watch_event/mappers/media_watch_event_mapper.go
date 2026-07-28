package mappers

import (
	"time"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
)

func (m *Mappers) MapMediaWatchEventDomainToEntity(event *ddownload.MediaWatchEvent) (*ewatchevent.MediaWatchEvent, error) {
	return &ewatchevent.MediaWatchEvent{
		EventID:    event.EventID,
		DownloadID: event.DownloadID,
		UserID:     event.UserID,
		SessionID:  event.SessionID,
		PositionMs: int(event.Position.Milliseconds()),
		IntervalMs: int(event.Interval.Milliseconds()),
	}, nil
}

func (m *Mappers) MapMediaWatchEventEntityToDomain(event *ewatchevent.MediaWatchEvent) (*ddownload.MediaWatchEvent, error) {
	return &ddownload.MediaWatchEvent{
		EventID:    event.EventID,
		DownloadID: event.DownloadID,
		UserID:     event.UserID,
		SessionID:  event.SessionID,
		Position:   time.Duration(event.PositionMs) * time.Millisecond,
		Interval:   time.Duration(event.IntervalMs) * time.Millisecond,
		CreatedAt:  event.CreatedAt,
	}, nil
}
