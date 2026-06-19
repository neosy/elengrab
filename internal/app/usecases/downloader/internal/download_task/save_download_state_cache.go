package dltask

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadTask) saveToDownloadStateCache(ctx context.Context, downloadID uuid.UUID, taskID uuid.UUID) {
	dlStateCache, _ := uc.dlStateCache.FindByDownloadID(ctx, downloadID)
	if dlStateCache != nil && dlStateCache.Download != nil {
		task, _ := uc.FindByTaskID(ctx, taskID)
		if task != nil {
			dlStateCache.TaskID = &task.TaskID
			dlStateCache.Download.DownloadTask = task
			err := uc.dlStateCache.Save(ctx, dlStateCache)
			if err != nil {
				uc.logger.Warn("Failed to save download state cache", "error", err)
			}
		}
	}
}
