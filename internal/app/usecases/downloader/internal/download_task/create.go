package dltask

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *DownloadTask) Create(ctx context.Context, task *ddownload.DownloadTask) error {
	if task == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	if task.TaskID == uuid.Nil {
		task.TaskID = uuid.New()
	}
	task.Status = dtypes.DownloadTaskStatusNew

	err := uc.TaskRep.Insert(ctx, task)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record",
			"error", err,
		)
		return err
	}

	uc.saveToDownloadStateCache(ctx, task.DownloadID, task.TaskID)

	return nil
}
