package dlstate

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadStateCache) DetachDownloadTask(ctx context.Context, taskID uuid.UUID) {
	dlStateCache, _ := uc.FindByTaskID(ctx, taskID)
	if dlStateCache != nil {
		dlStateCache.TaskID = nil
		if dlStateCache.Download != nil {
			dlStateCache.Download.DownloadTask = nil
		}
	}
}
