package uwatchposition

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaUserWatchPosition struct {
	logger *slog.Logger

	// repositories
	positionRepo persistence.MediaUserWatchPositionRepositoryFactory

	// in memory
	positionCacheRep persistence.MediaUserWatchPositionCacheRepository
}

func NewMediaUserWatchPosition(
	logger *slog.Logger,

	// repositories
	positionRepo persistence.MediaUserWatchPositionRepositoryFactory,

	// in memory
	positionCacheRep persistence.MediaUserWatchPositionCacheRepository,
) *MediaUserWatchPosition {
	return &MediaUserWatchPosition{
		logger:           logger,
		positionRepo:     positionRepo,
		positionCacheRep: positionCacheRep,
	}
}
