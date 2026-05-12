package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to done
func (s *MediaDownloadStatus) New(
	ctx context.Context,
	downloadID uuid.UUID,
) error {
	updateFieldsFunc := func(download *ddownload.MediaDownload) {
		download.ErrorMessage = nil
	}

	task, err := s.dlTask.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return err
	}

	err = s.dlTaskStatus.New(ctx, task.TaskID)
	if err != nil {
		return err
	}

	return s.updateStatus(
		ctx,
		downloadID,
		dtypes.MediaDownloadStatusNew,
		updateFieldsFunc,
	)
}
