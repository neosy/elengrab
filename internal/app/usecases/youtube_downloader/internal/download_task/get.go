package dltask

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *DownloadTask) FindByFileId(ctx context.Context, fileId uuid.UUID, checkNotFound bool) (*ddownload.DownloadTask, error) {
	task, err := uc.TaskRep.FindByFileId(ctx, fileId)
	if err != nil {
		uc.logger.Warn("Error finding record", "error", err)
		return nil, err
	}

	if checkNotFound && task == nil {
		uc.logger.Warn("Record not found", "fileId", fileId)
		return nil, err
	}

	return task, err
}

func (uc *DownloadTask) FindByTaskId(ctx context.Context, taskId uuid.UUID, checkNotFound bool) (*ddownload.DownloadTask, error) {
	task, err := uc.TaskRep.FindByTaskId(ctx, taskId)
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
