package mediadownload

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaDownload) Update(ctx context.Context, download *ddownload.MediaDownload) error {
	err := uc.downloadRep.Update(ctx, download)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}

	uc.saveToDownloadStateCache(ctx, download.DownloadID)

	return err
}
