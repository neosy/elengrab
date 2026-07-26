package mappers

import (
	"database/sql"
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

func (m *Mappers) MapMediaWatchEventRowsToDomain(rows *sql.Rows, fn func(*ddownload.MediaWatchEvent) error) error {
	var eEvent ewatchevent.MediaWatchEvent

	for rows.Next() {
		err := rows.Scan(eEvent.FieldPointers()...)
		if err != nil {
			return err
		}

		role, err := m.MapMediaWatchEventEntityToDomain(&eEvent)
		if err != nil {
			return err
		}

		err = fn(role)
		if err != nil {
			return err
		}
	}

	return nil
}
