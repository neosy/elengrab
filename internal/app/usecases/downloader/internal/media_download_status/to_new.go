package downloadstatus

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to done
func (s *mediaDownloadStatus) New(
	ctx context.Context,
	downloadID uuid.UUID,
) error {
	taskStatus, err := s.resolveTaskStatus()
	if err != nil {
		return err
	}

	mutate := func(download *ddownload.MediaDownload) error {
		download.ErrorMessage = nil
		return nil
	}

	taskID, err := s.getTaskID(ctx, downloadID)
	if err != nil {
		return err
	}

	err = taskStatus.New(ctx, taskID)
	if err != nil {
		return err
	}

	return s.updateStatus(
		ctx,
		downloadID,
		dtypes.MediaDownloadStatusNew,
		mutate,
	)
}
