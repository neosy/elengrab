package dltask

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *DownloadTask) Update(ctx context.Context, task *ddownload.DownloadTask) error {
	err := uc.taskRep.Update(ctx, task)
	if err != nil {
		uc.logger.Error("Update record error", "error", err)
		return err
	}
	return err
}
