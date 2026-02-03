package ddownload

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type DownloadTask struct {
	// Unique task identifier (UUID)
	TaskId uuid.UUID

	// Unique file identifier (UUID)
	FileId uuid.UUID

	// Status
	Status dtypes.DownloadTaskStatus

	// Youtube URL
	YoutubeUrl string

	// Youtube download options
	Options *DownloadOptions

	// Id worker
	WorkerId *uint

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
	copy.WorkerId = uptr.Copy(src.WorkerId)
	copy.JobID = uptr.Copy(src.JobID)

	return copy
}
