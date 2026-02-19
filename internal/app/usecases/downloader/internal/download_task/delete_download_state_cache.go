package dltask

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadTask) deleteToDownloadStateCache(ctx context.Context, taskID uuid.UUID) {
	dlStateCache, _ := uc.dlStateCache.FindByTaskID(ctx, taskID)
	if dlStateCache != nil {
		dlStateCache.TaskID = nil
		if dlStateCache.File != nil {
			dlStateCache.File.DownloadTask = nil
		}
	}
}
