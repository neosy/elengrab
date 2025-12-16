package dltask

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadTask) deleteToDownloadStateCache(ctx context.Context, taskID uuid.UUID) {
	dlStateCache, _ := uc.dlStateCache.FindByTaskId(ctx, taskID)
	if dlStateCache != nil {
		dlStateCache.TaskId = nil
		if dlStateCache.File != nil {
			dlStateCache.File.DownloadTask = nil
		}
	}
}
