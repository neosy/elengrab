package downloadstatus

import (
	"context"

	"uuid"

	itask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Pending set status to working
func (s *mediaDownloadStatus) Working(
	ctx context.Context,
	downloadID uuid.UUID,
	taskId uuid.UUID,
	workerId uint64,
) error {
	taskStatus, err := itask.SelectTaskStatus(s.sourceType, s.dlTaskStatus, s.processingTaskStatus)
	if err != nil {
		return err
	}

	err = taskStatus.Working(ctx, taskId, workerId)
	if err != nil {
		return err
	}

	return s.updateStatus(
		ctx,
		downloadID,
		dtypes.MediaDownloadStatusWorking,
		nil,
	)
}
