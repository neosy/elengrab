package downloader

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *YouTubeDownloader) ResetStuckJobs(ctx context.Context) error {
	err := uc.file.Tx(
		ctx,
		func(ctx context.Context) error {
			if err := uc.dlTask.ResetStatus(ctx); err != nil {
				return err
			}
			if err := uc.file.ResetStatus(ctx); err != nil {
				return err
			}
			if err := uc.file.DeleteBroken(ctx); err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	files, err := uc.file.GetByStatus(ctx, dtypes.FileStatusNew)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.DownloadTask != nil {
			uc.addFileToQueueDownload(ctx, file.FileID, file.DownloadTask.TaskID)
		}
	}

	return nil
}
