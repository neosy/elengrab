package mediadownload

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaDownload) SoftDelete(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.downloadRep.Delete(ctx, downloadID, true)
	if err != nil {
		uc.logger.Warn("Failed delete file", "error", err)
		return err
	}
	err = uc.dlStateCache.Delete(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed delete download state cache", "error", err)
	}
	return nil
}

func (uc *MediaDownload) HardDelete(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.downloadRep.Delete(ctx, downloadID, false)
	if err != nil {
		uc.logger.Warn("Failed delete file", "error", err)
		return err
	}
	err = uc.dlStateCache.Delete(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed delete download state cache", "error", err)
	}
	return nil
}
