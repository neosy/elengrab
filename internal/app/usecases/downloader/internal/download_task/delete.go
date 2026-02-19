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

func (uc *DownloadTask) DeleteByFileID(ctx context.Context, fileID uuid.UUID) error {
	return uc.TaskRep.DeleteByFileID(ctx, fileID)
}
