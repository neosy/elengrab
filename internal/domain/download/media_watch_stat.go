package ddownload

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type MediaWatchStat struct {
	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID

	// Number of completed views
	Views uint32

	// Record update timestamp, set automatically
	UpdatedAt time.Time
}

func (c *MediaWatchStat) Validate() error {
	if c.DownloadID == uuid.Nil {
		return errors.New("download ID is required")
	}

	if c.Views == 0 {
		return errors.New("views must be greater than zero")
	}

	return nil
}

func (src *MediaWatchStat) Copy() *MediaWatchStat {
	if src == nil {
		return nil
	}

	copy := *src

	return &copy
}
