package mediadownload

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *MediaDownload) ResetStatus(ctx context.Context) error {
	statuses := []dtypes.MediaDownloadStatus{
		dtypes.MediaDownloadStatusPending,
		dtypes.MediaDownloadStatusWorking,
	}

	err := uc.downloadRep.UpdateStatusToNew(ctx, statuses)
	if err != nil {
		uc.logger.Warn("Failed update status to new", "error", err)
		return err
	}

	return nil
}
