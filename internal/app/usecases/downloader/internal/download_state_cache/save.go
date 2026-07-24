package dlstate

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *DownloadStateCache) Save(ctx context.Context, state *ddownload.DownloadState) error {
	return uc.stateCacheRep.Save(ctx, state)
}

func (uc *DownloadStateCache) SaveByDownload(ctx context.Context, download *ddownload.MediaDownload) error {
	if download == nil {
		return nil
	}

	var taskId *uuid.UUID
	if download.DownloadTask != nil {
		taskId = &download.DownloadTask.TaskID
	}

	state := &ddownload.DownloadState{
		DownloadID: download.DownloadID,
		TaskID:     taskId,
		Download:   download,
	}

	return uc.stateCacheRep.Save(ctx, state)
}
