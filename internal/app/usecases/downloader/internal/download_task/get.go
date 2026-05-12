package dltask

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *DownloadTask) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.DownloadTask, error) {
	task, err := uc.TaskRep.FindByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Error finding record", "error", err)
		return nil, err
	}

	return task, err
}

func (uc *DownloadTask) GetByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.DownloadTask, error) {
	task, err := uc.FindByDownloadID(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		uc.logger.Warn("Record not found", "downloadID", downloadID)
		return nil, errorx.Errorf("task not found for downloadID: %s", downloadID, exceptionx.NOT_FOUND)
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
		return nil, errorx.Errorf("task not found for taskId: %s", taskId, exceptionx.NOT_FOUND)
	}

	return task, err
}
