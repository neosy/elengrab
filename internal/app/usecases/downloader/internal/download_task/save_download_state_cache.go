package dltask

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadTask) saveToDownloadStateCache(ctx context.Context, fileID uuid.UUID, taskID uuid.UUID) {
	dlStateCache, _ := uc.dlStateCache.FindByFileID(ctx, nil, fileID)
	if dlStateCache != nil && dlStateCache.File != nil {
		task, _ := uc.FindByTaskID(ctx, taskID, false)
		if task != nil {
			dlStateCache.TaskID = &task.TaskID
			dlStateCache.File.DownloadTask = task
			err := uc.dlStateCache.Save(ctx, dlStateCache)
			if err != nil {
				uc.logger.Warn("Failed to save download state cache", "error", err)
			}
		}
	}
}
