package mediadownload

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaDownload) Restore(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.downloadRep.Restore(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed restore file", "error", err)
		return err
	}
	return nil
}
