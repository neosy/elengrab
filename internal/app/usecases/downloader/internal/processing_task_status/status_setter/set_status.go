package statussetter

import (
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// SetStatus checks if the status transition is valid and setting the status.
func (u *StatusSetter) SetStatus(task *ddownload.MediaProcessingTask, toStatus dtypes.MediaTaskStatus) error {
	if task == nil {
		return apperrors.ErrFuncParamNullPointer
	}

	err := u.checkSetStatus(task.Status, toStatus)
	if err != nil {
		return err
	}

	// Update fields depending on the status being set
	switch toStatus {
	case dtypes.MediaTaskStatusWorking:
	}

	task.Status = toStatus

	return nil
}
