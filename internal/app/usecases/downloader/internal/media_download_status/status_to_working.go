package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Pending set status to working
func (uc *MediaDownloadStatus) Working(
	ctx context.Context,
	downloadID uuid.UUID,
	taskId uuid.UUID,
	workerId uint64,
) error {
	err := uc.dlTaskStatus.Working(ctx, taskId, workerId)
	if err != nil {
		return err
	}

	return uc.updateStatus(
		ctx,
		downloadID,
		dtypes.MediaDownloadStatusWorking,
		nil,
	)
}
