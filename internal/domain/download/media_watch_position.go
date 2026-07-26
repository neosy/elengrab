package ddownload

import (
	"errors"
	"time"

	"github.com/google/uuid"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type MediaWatchPosition struct {
	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID

	// Associated user identifier (UUID)
	UserID uuid.UUID

	// User session identifier (UUID)
	SessionID *uuid.UUID

	// Last saved playback position in milliseconds
	Position time.Duration

	// Record creation timestamp, set automatically
	CreatedAt time.Time

	// Record last update timestamp, set automatically
	UpdatedAt time.Time
}

func (p *MediaWatchPosition) Validate() error {
	if p.DownloadID == uuid.Nil {
		return errors.New("download ID is required")
	}

	if p.UserID == uuid.Nil && (p.SessionID == nil || *p.SessionID == uuid.Nil) {
		return errors.New("user ID is required")
	}

	return nil
}

func (src *MediaWatchPosition) Copy() *MediaWatchPosition {
	if src == nil {
		return nil
	}

	copy := new(*src)

	copy.SessionID = uptr.Copy(src.SessionID)

	return copy
}

func (p *MediaWatchPosition) Normalize(mediaDuration time.Duration) {
	if p.Position < 0 {
		p.Position = 0
		return
	}

	if mediaDuration-p.Position < time.Second {
		p.Position = 0
	}
}
