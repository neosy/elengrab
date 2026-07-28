package ddownload

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type MediaUserWatchChunk struct {
	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID

	// Associated user identifier (UUID)
	UserID uuid.UUID

	// Zero-based index of the 1000ms media chunk
	ChunkIndex uint32

	// How many times this chunk was watched
	Qty uint32

	// Record creation timestamp, set automatically
	CreatedAt time.Time
}

func (c *MediaUserWatchChunk) Validate() error {
	if c.DownloadID == uuid.Nil {
		return errors.New("download ID is required")
	}

	if c.Qty == 0 {
		return errors.New("quantity must be greater than zero")
	}

	return nil
}
