package mediadownload

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaDownload) saveUploadStateByDownload(ctx context.Context, download *ddownload.MediaDownload) error {
	if download == nil {
		return nil
	}

	state := ddownload.NewUploadState(download)

	oldState, _ := uc.uploadStateCacheRepo.FindByDownloadID(ctx, download.DownloadID)
	if oldState != nil && state.Download.MediaInfo == nil && oldState.Download.MediaInfo != nil {
		state.Download.MediaInfo = oldState.Download.MediaInfo
	}

	err := uc.uploadStateCacheRepo.Save(ctx, state)
	if err != nil {
		uc.logger.Warn("Failed to save upload state cache", "error", err)
		return err
	}

	return nil
}
