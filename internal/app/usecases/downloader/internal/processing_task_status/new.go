package proctaskstatus

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// New set status to new
func (s *ProcessingTaskStatus) New(
	ctx context.Context,
	taskId uuid.UUID,
) error {
	patchFields := func(task *ddownload.MediaProcessingTask) {
		task.WorkerID = nil
		task.JobID = nil
	}

	return s.updateStatus(
		ctx,
		taskId,
		dtypes.MediaTaskStatusNew,
		patchFields,
	)
}
