package dltaskstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Working set status to working
func (s *DownloadTaskStatus) Working(
	ctx context.Context,
	taskId uuid.UUID,
	workerId uint64,
) error {
	updateFieldsFunc := func(task *ddownload.DownloadTask) {
		task.WorkerID = &workerId
	}

	return s.updateStatus(
		ctx,
		taskId,
		dtypes.DownloadTaskStatusWorking,
		updateFieldsFunc,
	)
}
