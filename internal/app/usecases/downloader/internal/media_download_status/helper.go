package downloadstatus

import (
	"context"
	"errors"
	"uuid"

	itask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *mediaDownloadStatus) getTaskID(ctx context.Context, downloadID uuid.UUID) (uuid.UUID, error) {
	switch uc.sourceType {
	case dtypes.MediaDownloadSourceTypeURL:
		task, err := uc.dlTask.GetByDownloadID(ctx, downloadID)
		if err != nil {
			return uuid.Nil(), err
		}
		return task.TaskID, nil
	case dtypes.MediaDownloadSourceTypeUpload:
		task, err := uc.processingTask.GetByDownloadID(ctx, downloadID)
		if err != nil {
			return uuid.Nil(), err
		}
		return task.TaskID, nil
	default:
		err := errors.New("unsupported download source type")
		return uuid.Nil(), err
	}
}

func (uc *mediaDownloadStatus) resolveTaskStatus() (itask.TaskStatus, error) {
	return itask.TaskStatusBySourceType(uc.sourceType, uc.dlTaskStatus, uc.processingTaskStatus)
}
