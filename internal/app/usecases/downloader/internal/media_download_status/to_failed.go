package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	itask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

// Failed set status to failed
func (s *mediaDownloadStatus) Failed(
	ctx context.Context,
	downloadID uuid.UUID,
	mutate func(*ddownload.MediaDownload) error,
	message string,
) error {
	task, err := s.dlTask.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return err
	}

	taskStatus, err := itask.SelectTaskStatus(s.sourceType, s.dlTaskStatus, s.processingTaskStatus)
	if err != nil {
		return err
	}

	err = taskStatus.Failed(ctx, task.TaskID)
	if err != nil {
		return err
	}

	err = s.updateStatus(
		ctx,
		downloadID,
		dtypes.MediaDownloadStatusFailed,
		func(download *ddownload.MediaDownload) error {
			if mutate != nil {
				if err := mutate(download); err != nil {
					return err
				}
			}
			download.ErrorMessage = uptr.NonZeroString(message)
			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}
