package dltaskstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to failed
func (s *DownloadTaskStatus) Failed(
	ctx context.Context,
	taskId uuid.UUID,
) error {
	updateFieldsFunc := func(task *ddownload.DownloadTask) {
		task.WorkerId = nil
		task.JobID = nil
	}

	return s.updateStatus(
		ctx,
		taskId,
		dtypes.DownloadTaskStatusFailed,
		updateFieldsFunc,
	)
}
