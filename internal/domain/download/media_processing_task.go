package ddownload

import (
	"errors"
	"time"

	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type MediaProcessingTask struct {
	// Unique task identifier (UUID)
	TaskID uuid.UUID

	// Unique file identifier (UUID)
	DownloadID uuid.UUID

	// Status
	Status dtypes.MediaTaskStatus

	// Path to the staged uploaded file
	SourcePath string

	// Embeds TaskAssignment struct to track which worker and job are processing the task
	TaskAssignment

	// Task creation timestamp
	CreatedAt time.Time

	// Last update timestamp
	UpdatedAt time.Time
}

// Clone creates a deep copy of the MediaProcessingTask.
func (task *MediaProcessingTask) Clone() *MediaProcessingTask {
	if task == nil {
		return nil
	}

	copy := *task

	copy.WorkerID = uptr.Clone(task.WorkerID)
	copy.JobID = uptr.Clone(task.JobID)

	return &copy
}

// Validate checks the integrity of the MediaProcessingTask fields.
func (task *MediaProcessingTask) Validate() error {
	if task == nil {
		return errors.New("download task is nil")
	}

	if task.SourcePath == "" {
		return errors.New("source path is empty")
	}

	return nil
}
