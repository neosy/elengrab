package dlstate

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *DownloadStateCache) Save(ctx context.Context, state *ddownload.DownloadState) error {
	return uc.stateRep.Save(ctx, state)
}

func (uc *DownloadStateCache) SaveByFile(ctx context.Context, file *ddownload.File) error {
	if file == nil {
		return nil
	}

	var taskId *uuid.UUID
	if file.DownloadTask != nil {
		taskId = &file.DownloadTask.TaskID
	}

	state := &ddownload.DownloadState{
		FileID: file.FileID,
		TaskID: taskId,
		File:   file,
	}

	return uc.stateRep.Save(ctx, state)
}
