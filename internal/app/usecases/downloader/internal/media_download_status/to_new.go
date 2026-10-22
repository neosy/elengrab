package downloadstatus

import (
	"context"

	"uuid"

	itask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to done
func (s *mediaDownloadStatus) New(
	ctx context.Context,
	downloadID uuid.UUID,
) error {
	taskStatus, err := itask.SelectTaskStatus(s.sourceType, s.dlTaskStatus, s.processingTaskStatus)
	if err != nil {
		return err
	}

	mutate := func(download *ddownload.MediaDownload) error {
		download.ErrorMessage = nil
		return nil
	}

	task, err := s.dlTask.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return err
	}

	err = taskStatus.New(ctx, task.TaskID)
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
