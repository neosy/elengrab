package dltask

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadTask) saveToDownloadStateCache(ctx context.Context, fileID uuid.UUID, taskID uuid.UUID) {
	dlStateCache, _ := uc.dlStateCache.FindByFileId(ctx, nil, fileID)
	if dlStateCache != nil && dlStateCache.File != nil {
		task, _ := uc.FindByTaskId(ctx, taskID, false)
		if task != nil {
			dlStateCache.TaskId = &task.TaskId
			dlStateCache.File.DownloadTask = task
			err := uc.dlStateCache.Save(ctx, dlStateCache)
			if err != nil {
				uc.logger.Warn("Failed to save download state cache", "error", err)
			}
		}
	}
}
