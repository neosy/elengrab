package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to failed
func (s *MediaDownloadStatus) Failed(
	ctx context.Context,
	downloadID uuid.UUID,
	patch *dto.MediaDownloadInfoPatch,
	message *string,
) error {
	updateFieldsFunc := func(download *ddownload.MediaDownload) {
		dto.PatchToMediaDownloadDomain(patch, download)
		download.ErrorMessage = message
	}

	task, err := s.dlTask.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return err
	}

	err = s.dlTaskStatus.Failed(ctx, task.TaskID)
	if err != nil {
		return err
	}

	return s.updateStatus(
		ctx,
		downloadID,
		dtypes.MediaDownloadStatusFailed,
		updateFieldsFunc,
	)
}
