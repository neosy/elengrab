package uwatchposition

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaUserWatchPosition struct {
	logger *slog.Logger

	// repositories
	positionRep persistence.MediaUserWatchPositionRepository

	// in memory
	positionCacheRep persistence.MediaUserWatchPositionCacheRepository
}

func NewMediaUserWatchPosition(
	logger *slog.Logger,

	// repositories
	positionRep persistence.MediaUserWatchPositionRepository,

	// in memory
	positionCacheRep persistence.MediaUserWatchPositionCacheRepository,
) *MediaUserWatchPosition {
	return &MediaUserWatchPosition{
		logger:           logger,
		positionRep:      positionRep,
		positionCacheRep: positionCacheRep,
	}
}
