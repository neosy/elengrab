package ddownload

import (
	"errors"
	"time"

	"github.com/google/uuid"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type MediaWatchEvent struct {
	// Unique event identifier (UUID)
	EventID uuid.UUID

	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID

	// Associated user identifier (UUID)
	UserID *uuid.UUID

	// User session identifier (UUID)
	SessionID *uuid.UUID

	// Playback position
	Position time.Duration

	// Playback duration since the previous event
	Interval time.Duration

	// Record creation timestamp, set automatically
	CreatedAt time.Time
}

func (e *MediaWatchEvent) Start() time.Duration {
	return max(e.Position-e.Interval, 0)
}

func (e *MediaWatchEvent) Validate() error {
	if e.DownloadID == uuid.Nil {
		return errors.New("download ID is required")
	}

	if e.Position < 0 {
		return errors.New("position must not be negative")
	}

	if e.Interval <= 0 {
		return errors.New("interval must be greater than zero")
	}

	return nil
}

func (e *MediaWatchEvent) Copy() *MediaWatchEvent {
	if e == nil {
		return nil
	}

	copy := new(*e)

	copy.UserID = uptr.Copy(e.UserID)
	copy.SessionID = uptr.Copy(e.SessionID)

	return copy
}
