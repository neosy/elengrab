package proctaskstatus

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// updateStatus marking the status
func (s *ProcessingTaskStatus) updateStatus(
	ctx context.Context,
	taskId uuid.UUID,
	toStatus dtypes.MediaTaskStatus,
	patchFields func(file *ddownload.MediaProcessingTask),
) error {
	task, err := s.processingTask.GetByTaskID(ctx, taskId)
	if err != nil {
		return err
	}

	err = s.statusSetter.SetStatus(task, toStatus)
	if err != nil {
		s.logger.Warn(
			"Failed to update status",
			"taskId", taskId,
			"error", err,
		)
		return err
	}

	// Patch fields
	if patchFields != nil {
		patchFields(task)
	}

	// Update in the repository
	err = s.processingTask.Update(ctx, task)
	if err != nil {
		s.logger.Warn(
			"Failed to update task in the repository",
			"taskId", taskId,
			"error", err,
		)
		return err
	}

	return nil
}
