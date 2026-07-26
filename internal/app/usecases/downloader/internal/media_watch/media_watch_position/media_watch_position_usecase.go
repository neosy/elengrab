package watchposition

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaWatchPosition struct {
	logger *slog.Logger

	// repositories
	positionRep persistence.MediaWatchPositionRepository

	// in memory
	positionCacheRep persistence.MediaWatchPositionCacheRepository
}

func NewMediaWatchPosition(
	logger *slog.Logger,

	// repositories
	positionRep persistence.MediaWatchPositionRepository,

	// in memory
	positionCacheRep persistence.MediaWatchPositionCacheRepository,
) *MediaWatchPosition {
	return &MediaWatchPosition{
		logger:           logger,
		positionRep:      positionRep,
		positionCacheRep: positionCacheRep,
	}
}
