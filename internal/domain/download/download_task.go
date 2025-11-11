package ddownload

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type DownloadTask struct {
	// Unique task identifier (UUID)
	TaskId uuid.UUID

	// Unique file identifier (UUID)
	FileId uuid.UUID

	// Status
	Status dtypes.DownloadTaskStatus

	// Id worker
	WorkerId *uint

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}
