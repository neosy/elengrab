package dltask

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *DownloadTask) Update(ctx context.Context, task *ddownload.DownloadTask) error {
	err := uc.TaskRep.Update(ctx, task)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}

	uc.saveToDownloadStateCache(ctx, task.FileID, task.TaskID)

	return err
}
