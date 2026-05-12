package downloader

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *Downloader) ResetStuckJobs(ctx context.Context) error {
	err := uc.download.Tx(
		ctx,
		func(ctx context.Context) error {
			if err := uc.dlTask.ResetStatus(ctx); err != nil {
				return err
			}
			if err := uc.download.ResetStatus(ctx); err != nil {
				return err
			}
			if err := uc.download.DeleteBroken(ctx); err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	downloads, err := uc.download.GetByStatus(ctx, dtypes.MediaDownloadStatusNew)
	if err != nil {
		return err
	}

	for _, download := range downloads {
		if download.DownloadTask != nil {
			uc.addDownloadToQueueDownload(ctx, download.DownloadID, download.DownloadTask.TaskID)
		}
	}

	return nil
}
