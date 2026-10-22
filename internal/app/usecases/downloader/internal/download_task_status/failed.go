package dltaskstatus

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to failed
func (s *DownloadTaskStatus) Failed(
	ctx context.Context,
	taskId uuid.UUID,
) error {
	patchFields := func(task *ddownload.DownloadTask) {
		task.WorkerID = nil
		task.JobID = nil
	}

	return s.updateStatus(
		ctx,
		taskId,
		dtypes.MediaTaskStatusFailed,
		patchFields,
	)
}
