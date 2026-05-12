package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to failed
func (s *MediaDownloadStatus) Failed(
	ctx context.Context,
	downloadID uuid.UUID,
	patch func(*ddownload.MediaDownload),
	message *string,
) error {
	task, err := s.dlTask.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return err
	}

	err = s.dlTaskStatus.Failed(ctx, task.TaskID)
	if err != nil {
		return err
	}

	err = s.updateStatus(
		ctx,
		downloadID,
		dtypes.MediaDownloadStatusFailed,
		func(download *ddownload.MediaDownload) {
			patch(download)
			download.ErrorMessage = message
		},
	)
	if err != nil {
		return err
	}

	return nil
}
