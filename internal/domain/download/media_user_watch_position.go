package ddownload

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type MediaUserWatchPosition struct {
	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID

	// Associated user identifier (UUID)
	UserID uuid.UUID

	// User session identifier (UUID)
	SessionID uuid.UUID

	// Last saved playback position in milliseconds
	Position time.Duration

	// Record creation timestamp, set automatically
	CreatedAt time.Time

	// Record last update timestamp, set automatically
	UpdatedAt time.Time
}

func (p *MediaUserWatchPosition) Validate() error {
	if p.DownloadID == uuid.Nil {
		return errors.New("download ID is required")
	}

	if p.UserID == uuid.Nil && p.SessionID == uuid.Nil {
		return errors.New("user ID is required")
	}

	return nil
}

func (src *MediaUserWatchPosition) Copy() *MediaUserWatchPosition {
	if src == nil {
		return nil
	}

	copy := new(*src)

	return copy
}

func (p *MediaUserWatchPosition) Normalize(mediaDuration time.Duration) {
	if p.Position < 0 {
		p.Position = 0
		return
	}

	if mediaDuration-p.Position < time.Second {
		p.Position = 0
	}
}
