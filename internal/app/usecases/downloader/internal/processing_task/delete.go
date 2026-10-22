package processingtask

import (
	"context"

	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *ProcessingTask) Delete(ctx context.Context, taskID uuid.UUID) error {
	err := uc.TaskRepo().Delete(ctx, taskID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ProcessingTask) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	return uc.TaskRepo().DeleteByDownloadID(ctx, downloadID)
}

func (uc *ProcessingTask) DeleteByStatus(ctx context.Context, status dtypes.MediaTaskStatus) error {
	return uc.TaskRepo().DeleteByStatus(ctx, status)
}
