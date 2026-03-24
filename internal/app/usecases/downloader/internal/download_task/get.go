package dltask

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

func (uc *DownloadTask) FindByFileID(ctx context.Context, fileId uuid.UUID) (*ddownload.DownloadTask, error) {
	task, err := uc.TaskRep.FindByFileID(ctx, fileId)
	if err != nil {
		uc.logger.Warn("Error finding record", "error", err)
		return nil, err
	}

	return task, err
}

func (uc *DownloadTask) GetByFileID(ctx context.Context, fileId uuid.UUID) (*ddownload.DownloadTask, error) {
	task, err := uc.FindByFileID(ctx, fileId)
	if err != nil {
		return nil, err
	}

	if task == nil {
		uc.logger.Warn("Record not found", "fileId", fileId)
		return nil, errorx.NewFromError(
			fmt.Errorf("task not found for fileId: %s", fileId),
			exceptionx.NOT_FOUND,
		)
	}

	return task, err
}

func (uc *DownloadTask) FindByTaskID(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadTask, error) {
	task, err := uc.TaskRep.FindByTaskID(ctx, taskId)
	if err != nil {
		uc.logger.Warn("Error finding record", "error", err)
		return nil, err
	}

	return task, err
}

func (uc *DownloadTask) GetByTaskID(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadTask, error) {
	task, err := uc.FindByTaskID(ctx, taskId)
	if err != nil {
		return nil, err
	}

	if task == nil {
		uc.logger.Warn("Record not found", "taskId", taskId)
		return nil, errorx.NewFromError(
			fmt.Errorf("task not found for taskId: %s", taskId),
			exceptionx.NOT_FOUND,
		)
	}

	return task, err
}
