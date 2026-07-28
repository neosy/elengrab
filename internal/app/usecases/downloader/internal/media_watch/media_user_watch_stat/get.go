package uwatchstat

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaUserWatchStat) Find(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) (*ddownload.MediaUserWatchStat, error) {
	if downloadID == uuid.Nil {
		return nil, nil
	}

	stat, err := uc.statRep.Find(ctx, downloadID, userID)
	if err != nil {
		uc.logger.Warn(
			"Failed to find media user watch statistics",
			"downloadID", downloadID,
			"userID", userID,
			"error", err,
		)
		return nil, errorx.Errorf("failed to find media user watch statistics: %w", err, exceptionx.ERROR)
	}

	return stat, nil
}
