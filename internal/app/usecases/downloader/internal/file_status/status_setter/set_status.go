package statussetter

import (
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// SetStatus checks if the status transition is valid and setting the status.
func (u *FileStatusSetter) SetStatus(file *ddownload.File, toStatus dtypes.FileStatus) error {
	if file == nil {
		return apperrors.ErrFuncParamNullPointer
	}

	err := u.checkSetStatus(file.Status, toStatus)
	if err != nil {
		return err
	}

	// Update fields depending on the status being set
	switch toStatus {
	case dtypes.FileStatusFailed:
	}

	file.Status = toStatus

	return nil
}
