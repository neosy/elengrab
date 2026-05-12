package dltask

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadTask) Delete(ctx context.Context, taskID uuid.UUID) error {
	err := uc.TaskRep.Delete(ctx, taskID)

	if err == nil {
		uc.deleteToDownloadStateCache(ctx, taskID)
	}

	return err
}

func (uc *DownloadTask) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	return uc.TaskRep.DeleteByDownloadID(ctx, downloadID)
}
