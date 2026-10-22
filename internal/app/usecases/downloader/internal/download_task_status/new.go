package dltaskstatus

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// New set status to new
func (s *DownloadTaskStatus) New(
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
		dtypes.MediaTaskStatusNew,
		patchFields,
	)
}
