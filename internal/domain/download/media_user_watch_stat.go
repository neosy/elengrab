package ddownload

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type MediaUserWatchStat struct {
	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID

	// Associated user identifier (UUID)
	// Use '00000000-0000-0000-0000-000000000000' for anonymous users.
	UserID uuid.UUID

	// Number of completed views
	Views uint32

	// Record update timestamp, set automatically
	UpdatedAt time.Time
}

func (c *MediaUserWatchStat) Validate() error {
	if c.DownloadID == uuid.Nil {
		return errors.New("download ID is required")
	}

	if c.Views == 0 {
		return errors.New("views must be greater than zero")
	}

	return nil
}

func (src *MediaUserWatchStat) Copy() *MediaUserWatchStat {
	if src == nil {
		return nil
	}

	copy := *src

	return &copy
}
