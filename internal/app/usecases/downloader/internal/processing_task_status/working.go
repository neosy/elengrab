package proctaskstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Working set status to working
func (s *ProcessingTaskStatus) Working(
	ctx context.Context,
	taskId uuid.UUID,
	workerId uint64,
) error {
	patchFields := func(task *ddownload.MediaProcessingTask) {
		task.WorkerID = &workerId
	}

	return s.updateStatus(
		ctx,
		taskId,
		dtypes.MediaTaskStatusWorking,
		patchFields,
	)
}
