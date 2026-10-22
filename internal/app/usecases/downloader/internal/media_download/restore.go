package mediadownload

import (
	"context"

	"uuid"
)

func (uc *MediaDownload) Restore(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.downloadRepo().Restore(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed restore file", "error", err)
		return err
	}

	uc.downloadCacheRep.Delete(ctx, downloadID)

	return nil
}
