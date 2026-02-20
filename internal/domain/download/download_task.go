package ddownload

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type DownloadTask struct {
	// Unique task identifier (UUID)
	TaskID uuid.UUID

	// Unique file identifier (UUID)
	FileID uuid.UUID

	// Status
	Status dtypes.DownloadTaskStatus

	// Media URL
	MediaUrl string

	// Media download options
	Options *DownloadOptions

	// Id worker
	WorkerID *uint64

	// ID job
	JobID *uuid.UUID

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}

func (src *DownloadTask) Copy() *DownloadTask {
	if src == nil {
		return nil
	}

	copy := uptr.Copy(src)
	copy.Options = src.Options.Copy()
	copy.WorkerID = uptr.Copy(src.WorkerID)
	copy.JobID = uptr.Copy(src.JobID)

	return copy
}
