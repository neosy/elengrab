package dlstate

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadStateCache) Delete(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.stateCacheRep.Delete(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed delete download state cache", "error", err)
		return err
	}

	return nil
}
