package uwatchstat

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaUserWatchStat) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.statRepo().DeleteByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Warn(
			"Failed to delete media user watch statistics",
			"downloadID", downloadID,
			"error", err,
		)
		return err
	}

	return nil
}

func (uc *MediaUserWatchStat) DeleteAll(ctx context.Context) error {
	err := uc.statRepo().DeleteAll(ctx)
	if err != nil {
		uc.logger.Warn(
			"Failed to delete all media user watch statistics",
			"error", err,
		)
		return err
	}

	return nil
}
