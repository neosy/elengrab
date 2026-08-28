package watchstat

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatchStat) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.statRepo().DeleteByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Warn(
			"Failed to delete media watch statistics",
			"downloadID", downloadID,
			"error", err,
		)
		return err
	}

	uc.statCacheRep.Delete(ctx, downloadID)

	return nil
}

func (uc *MediaWatchStat) DeleteAll(ctx context.Context) error {
	err := uc.statRepo().DeleteAll(ctx)
	if err != nil {
		uc.logger.Warn(
			"Failed to delete all media watch statistics",
			"error", err,
		)
		return err
	}

	return nil
}
