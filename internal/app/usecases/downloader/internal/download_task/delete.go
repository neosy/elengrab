package dltask

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadTask) Delete(ctx context.Context, taskID uuid.UUID) error {
	err := uc.TaskRepo().Delete(ctx, taskID)
	if err != nil {
		return err
	}

	uc.dlStateCache.DetachDownloadTask(ctx, taskID)

	return nil
}

func (uc *DownloadTask) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	return uc.TaskRepo().DeleteByDownloadID(ctx, downloadID)
}
