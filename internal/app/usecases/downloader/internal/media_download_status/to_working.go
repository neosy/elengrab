package downloadstatus

import (
	"context"

	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Pending set status to working
func (s *mediaDownloadStatus) Working(
	ctx context.Context,
	downloadID uuid.UUID,
	taskId uuid.UUID,
	workerId uint64,
) error {
	taskStatus, err := s.resolveTaskStatus()
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
