package downloader

import (
	"context"
	"fmt"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *downloader) ResetStuckJobs(ctx context.Context) error {
	err := uc.download.Tx(
		ctx,
		func(ctx context.Context) error {
			if err := uc.dlTask.ResetStatus(ctx); err != nil {
				return err
			}
			if err := uc.processingTask.ResetStatus(ctx); err != nil {
				return err
			}
			if err := uc.download.ResetStatus(ctx); err != nil {
				return err
			}
			if err := uc.download.DeleteNewWithoutTask(ctx); err != nil {
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
		switch download.SourceType {
		case dtypes.MediaDownloadSourceTypeURL:
			if download.DownloadTask != nil {
				uc.addDownloadToQueueDownload(ctx, download.DownloadID, download.DownloadTask.TaskID)
			}
		case dtypes.MediaDownloadSourceTypeUpload:
			if download.ProcessingTask != nil {
				uc.addUploadToQueueProcessing(ctx, download.DownloadID, download.ProcessingTask.TaskID)
			}
		default:
			return fmt.Errorf("unsupported source type %v for download %v", download.SourceType, download.DownloadID)
		}
	}

	return nil
}
