package watchstat

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaWatchStat) Find(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaWatchStat, error) {
	if downloadID == uuid.Nil {
		return nil, nil
	}

	stat, cacheStatus, _ := uc.statCacheRep.Find(ctx, downloadID)
	if stat != nil {
		return stat, nil
	}
	if cacheStatus == memsimple.CacheStatusNegativeHit {
		return nil, nil
	}

	stat, err := uc.statRep.Find(ctx, downloadID)
	if err != nil {
		uc.logger.Warn(
			"Failed to find media watch statistics",
			"downloadID", downloadID,
			"error", err,
		)
		return nil, errorx.Errorf("failed to find media watch statistics: %w", err, exceptionx.ERROR)
	}

	if stat != nil {
		uc.statCacheRep.Save(ctx, stat)
	} else {
		uc.statCacheRep.SaveNegative(ctx, downloadID)
	}

	return stat, nil
}

func (uc *MediaWatchStat) Exists(ctx context.Context, downloadID uuid.UUID) (bool, error) {
	stat, err := uc.Find(ctx, downloadID)
	return stat != nil, err
}
