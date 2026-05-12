package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// updateStatus marking the status
func (uc *MediaDownloadStatus) updateStatus(
	ctx context.Context,
	downloadID uuid.UUID,
	toStatus dtypes.MediaDownloadStatus,
	patch func(download *ddownload.MediaDownload),
) error {
	download, err := uc.download.GetByDownloadID(ctx, nil, downloadID)
	if err != nil {
		return err
	}

	err = uc.statusSetter.SetStatus(download, toStatus)
	if err != nil {
		uc.logger.Warn(
			"Failed to update status",
			"downloadID", downloadID,
			"error", err,
		)
		return err
	}

	// Update fields
	if patch != nil {
		patch(download)
	}

	// Update in the repository
	err = uc.download.Update(ctx, download)
	if err != nil {
		uc.logger.Warn(
			"Failed to update download in the repository",
			"downloadID", downloadID,
			"error", err,
		)

		return err
	}

	return nil
}
