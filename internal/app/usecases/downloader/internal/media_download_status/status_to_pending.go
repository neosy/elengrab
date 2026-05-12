package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Pending set status to pending
func (s *MediaDownloadStatus) Pending(
	ctx context.Context,
	downloadID uuid.UUID,
	taskId uuid.UUID,
	jobID uuid.UUID,
) error {
	err := s.dlTaskStatus.Penging(ctx, taskId, jobID)
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
