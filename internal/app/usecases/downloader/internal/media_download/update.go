package mediadownload

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaDownload) Update(ctx context.Context, userID *uuid.UUID, download *ddownload.MediaDownload) error {
	download.NormalizeForSave()

	if userID != nil && *userID != uuid.Nil && (download.UserID == nil || *download.UserID == uuid.Nil) {
		download.UserID = userID
	}

	err := uc.downloadRep.Update(ctx, download)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}

	uc.saveToDownloadStateCache(ctx, download.DownloadID)

	return err
}
