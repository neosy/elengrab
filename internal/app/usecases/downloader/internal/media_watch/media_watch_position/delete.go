package watchposition

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatchPosition) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	if downloadID == uuid.Nil {
		return nil
	}

	err := uc.positionRep.DeleteByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Warn(
			"Failed to delete media watch positions",
			"downloadID", downloadID,
			"error", err,
		)
		return err
	}

	return nil
}
