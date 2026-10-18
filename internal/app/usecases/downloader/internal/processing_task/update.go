package processingtask

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *ProcessingTask) Update(ctx context.Context, task *ddownload.MediaProcessingTask) error {
	err := uc.TaskRepo().Update(ctx, task)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}

	return err
}
