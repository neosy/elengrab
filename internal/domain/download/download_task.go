package ddownload

import (
	"errors"
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type DownloadTask struct {
	// Unique task identifier (UUID)
	TaskID uuid.UUID

	// Unique file identifier (UUID)
	DownloadID uuid.UUID

	// Status
	Status dtypes.MediaTaskStatus

	// Media URL
	MediaUrl string

	// Media download options
	Options *DownloadOptions

	// Embeds TaskAssignment struct to track which worker and job are processing the task
	TaskAssignment

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}

func (task *DownloadTask) Clone() *DownloadTask {
	if task == nil {
		return nil
	}

	copy := *task

	copy.Options = task.Options.Clone()
	copy.WorkerID = uptr.Clone(task.WorkerID)
	copy.JobID = uptr.Clone(task.JobID)

	return &copy
}

func (task *DownloadTask) Validate() error {
	if task == nil {
		return errors.New("download task is nil")
	}

	if task.MediaUrl == "" {
		return errors.New("media url is required")
	}

	return nil
}
