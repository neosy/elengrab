package dltaskstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// updateStatus marking the status
func (s *DownloadTaskStatus) updateStatus(
	ctx context.Context,
	taskId uuid.UUID,
	toStatus dtypes.DownloadTaskStatus,
	updateFieldsFunc func(file *ddownload.DownloadTask),
) error {
	task, err := s.downloadTask.FindByTaskId(ctx, taskId, true)
	if err != nil {
		return err
	}

	err = s.statusSetter.SetStatus(task, toStatus)
	if err != nil {
		s.logger.Debug(
			"Failed to update status",
			"taskId", taskId,
			"error", err,
		)
		return err
	}

	// Update fields
	if updateFieldsFunc != nil {
		updateFieldsFunc(task)
	}

	// Update in the repository
	err = s.downloadTask.Update(ctx, task)
	if err != nil {
		s.logger.Debug(
			"Failed to update file in the repository",
			"taskId", taskId,
			"error", err,
		)
		return err
	}

	return nil
}
