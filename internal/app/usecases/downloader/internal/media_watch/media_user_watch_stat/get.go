package uwatchstat

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaUserWatchStat) Find(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
) (*ddownload.MediaUserWatchStat, error) {
	if downloadID == uuid.Nil {
		return nil, nil
	}

	stat, cacheStatus, _ := uc.statCacheRep.Find(downloadID, userID)
	if stat != nil {
		return stat, nil
	}
	if cacheStatus == memsimple.CacheStatusNegativeHit {
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

	if stat != nil {
		uc.statCacheRep.Save(stat)
	} else {
		uc.statCacheRep.SaveNegative(downloadID, userID)
	}

	return stat, nil
}

func (uc *MediaUserWatchStat) Exists(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
) (bool, error) {
	stat, err := uc.Find(ctx, downloadID, userID)
	return stat != nil, err
}
