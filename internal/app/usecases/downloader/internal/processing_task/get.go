package processingtask

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *ProcessingTask) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaProcessingTask, error) {
	task, err := uc.TaskRepo().FindByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Error finding record", "error", err)
		return nil, err
	}

	return task, err
}

func (uc *ProcessingTask) GetByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaProcessingTask, error) {
	task, err := uc.FindByDownloadID(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		uc.logger.Warn("Record not found", "downloadID", downloadID)
		return nil, errorx.Errorf("processing task not found for downloadID: %s", downloadID, exceptionx.NOT_FOUND)
	}

	return task, err
}

func (uc *ProcessingTask) FindByTaskID(ctx context.Context, taskId uuid.UUID) (*ddownload.MediaProcessingTask, error) {
	task, err := uc.TaskRepo().FindByTaskID(ctx, taskId)
	if err != nil {
		uc.logger.Warn("Error finding record", "error", err)
		return nil, err
	}

	return task, err
}

func (uc *ProcessingTask) GetByTaskID(ctx context.Context, taskId uuid.UUID) (*ddownload.MediaProcessingTask, error) {
	task, err := uc.FindByTaskID(ctx, taskId)
	if err != nil {
		return nil, err
	}

	if task == nil {
		uc.logger.Warn("Record not found", "taskId", taskId)
		return nil, errorx.Errorf("processing task not found for taskId: %s", taskId, exceptionx.NOT_FOUND)
	}

	return task, err
}
