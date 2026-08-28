package mediadownload

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *MediaDownload) ResetStatus(ctx context.Context) error {
	statusesToNew := []dtypes.MediaDownloadStatus{
		dtypes.MediaDownloadStatusPending,
		dtypes.MediaDownloadStatusWorking,
	}

	err := uc.downloadRep().UpdateStatus(ctx, statusesToNew, dtypes.MediaDownloadStatusNew)
	if err != nil {
		uc.logger.Warn("Failed update status to new", "error", err)
		return err
	}

	statusesToDone := []dtypes.MediaDownloadStatus{
		dtypes.MediaDownloadStatusRefreshing,
	}

	err = uc.downloadRep().UpdateStatus(ctx, statusesToDone, dtypes.MediaDownloadStatusDone)
	if err != nil {
		uc.logger.Warn("Failed update status to done", "error", err)
		return err
	}

	return nil
}
