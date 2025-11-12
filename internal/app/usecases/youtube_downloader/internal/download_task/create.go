package dltask

import (
	"context"
	"errors"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *DownloadTask) Create(ctx context.Context, task *ddownload.DownloadTask) error {
	if task == nil {
		uc.logger.Error("Nil pointer in function")
		return errors.New("function parameter is a null pointer")
	}

	if task.TaskId == uuid.Nil {
		task.TaskId = uuid.New()
	}
	task.Status = dtypes.DownloadTaskStatusPending

	err := uc.taskRep.Insert(ctx, task)
	if err != nil {
		uc.logger.Error(
			"Failed to insert record into repository",
			"error", err,
		)
		return err
	}

	return nil
}
