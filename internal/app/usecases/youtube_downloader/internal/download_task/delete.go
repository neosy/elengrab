package dltask

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadTask) Delete(ctx context.Context, taskId uuid.UUID) error {
	return uc.TaskRep.Delete(ctx, taskId)
}

func (uc *DownloadTask) DeleteByFileId(ctx context.Context, fileId uuid.UUID) error {
	return uc.TaskRep.DeleteByFileId(ctx, fileId)
}
