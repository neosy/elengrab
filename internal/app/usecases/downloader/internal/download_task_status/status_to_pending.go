package dltaskstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Penging set status to penging
func (s *DownloadTaskStatus) Penging(
	ctx context.Context,
	taskId uuid.UUID,
	jobID uuid.UUID,
) error {
	updateFieldsFunc := func(task *ddownload.DownloadTask) {
		task.JobID = &jobID
	}

	return s.updateStatus(
		ctx,
		taskId,
		dtypes.DownloadTaskStatusPending,
		updateFieldsFunc,
	)
}
