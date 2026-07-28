package uwatchposition

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaUserWatchPosition) Find(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID,
) (*ddownload.MediaUserWatchPosition, error) {
	if downloadID == uuid.Nil {
		return nil, nil
	}

	if userID == uuid.Nil && (sessionID == nil || *sessionID == uuid.Nil) {
		return nil, nil
	}

	if userID != uuid.Nil && sessionID != nil {
		sessionID = nil
	}

	position, cacheStatus, _ := uc.positionCacheRep.Find(downloadID, userID, sessionID)
	if position != nil {
		return position, nil
	}
	if cacheStatus == memsimple.CacheStatusNegativeHit {
		return nil, nil
	}

	position, err := uc.positionRep.Find(ctx, downloadID, userID, sessionID)
	if err != nil {
		uc.logger.Warn(
			"Failed to find media watch positions",
			"downloadID", downloadID,
			"error", err,
		)
		return nil, errorx.Errorf("failed to find media watch positions: %w", err, exceptionx.ERROR)
	}

	if position != nil {
		uc.positionCacheRep.Save(position)
	} else {
		uc.positionCacheRep.SaveNegative(downloadID, userID, sessionID)
	}

	return position, nil
}
