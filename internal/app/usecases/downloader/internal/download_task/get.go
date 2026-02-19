package dltask

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

func (uc *DownloadTask) FindByFileID(ctx context.Context, fileId uuid.UUID, checkNotFound bool) (*ddownload.DownloadTask, error) {
	task, err := uc.TaskRep.FindByFileID(ctx, fileId)
	if err != nil {
		uc.logger.Warn("Error finding record", "error", err)
		return nil, err
	}

	if checkNotFound && task == nil {
		uc.logger.Warn("Record not found", "fileId", fileId)
		return nil, errorx.NewByErr(fmt.Errorf("task not found for fileId: %s", fileId), exceptionx.NOT_FOUND)
	}

	return task, err
}

func (uc *DownloadTask) FindByTaskID(ctx context.Context, taskId uuid.UUID, checkNotFound bool) (*ddownload.DownloadTask, error) {
	task, err := uc.TaskRep.FindByTaskID(ctx, taskId)
	if err != nil {
		uc.logger.Warn("Error finding record", "error", err)
		return nil, err
	}

	if checkNotFound && task == nil {
		uc.logger.Warn("Record not found", "taskId", taskId)
		return nil, err
	}

	return task, err
}
