package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	itask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Pending set status to pending
func (s *mediaDownloadStatus) Pending(
	ctx context.Context,
	downloadID uuid.UUID,
	taskId uuid.UUID,
	jobID uuid.UUID,
) error {
	taskStatus, err := itask.SelectTaskStatus(s.sourceType, s.dlTaskStatus, s.processingTaskStatus)
	if err != nil {
		return err
	}

	err = taskStatus.Penging(ctx, taskId, jobID)
	if err != nil {
		return err
	}

	return s.updateStatus(
		ctx,
		downloadID,
		dtypes.MediaDownloadStatusPending,
		nil,
	)
}
