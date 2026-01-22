package statussetter

import (
	"errors"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// SetStatus checks if the status transition is valid and setting the status.
func (u *DownloadTaskStatusSetter) SetStatus(task *ddownload.DownloadTask, toStatus dtypes.DownloadTaskStatus) error {
	if task == nil {
		return errors.New("function parameter is a null pointer")
	}

	err := u.checkSetStatus(task.Status, toStatus)
	if err != nil {
		return err
	}

	// Update fields depending on the status being set
	switch toStatus {
	case dtypes.DownloadTaskStatusWorking:
	}

	task.Status = toStatus

	return nil
}
