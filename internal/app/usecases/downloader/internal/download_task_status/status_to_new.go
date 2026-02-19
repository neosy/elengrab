package dltaskstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// New set status to new
func (s *DownloadTaskStatus) New(
	ctx context.Context,
	taskId uuid.UUID,
) error {
	updateFieldsFunc := func(task *ddownload.DownloadTask) {
		task.WorkerID = nil
		task.JobID = nil
	}

	return s.updateStatus(
		ctx,
		taskId,
		dtypes.DownloadTaskStatusNew,
		updateFieldsFunc,
	)
}
