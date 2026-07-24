package watchstat

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatchStat) Delete(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.statRep.Delete(ctx, downloadID)
	if err != nil {
		uc.logger.Warn(
			"Failed to delete media watch statistics",
			"downloadID", downloadID,
			"error", err,
		)
		return err
	}

	uc.statCacheRep.Delete(downloadID)

	return nil
}
