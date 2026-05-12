package statussetter

import (
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// SetStatus checks if the status transition is valid and setting the status.
func (u *MediaDownloadStatusSetter) SetStatus(download *ddownload.MediaDownload, toStatus dtypes.MediaDownloadStatus) error {
	if download == nil {
		return apperrors.ErrFuncParamNullPointer
	}

	err := u.checkSetStatus(download.Status, toStatus)
	if err != nil {
		return err
	}

	// Update fields depending on the status being set
	switch toStatus {
	case dtypes.MediaDownloadStatusFailed:
	}

	download.Status = toStatus

	return nil
}
