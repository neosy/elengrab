package uwatchstat

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaUserWatchStat struct {
	logger *slog.Logger

	// repositories
	statRepo persistence.MediaUserWatchStatRepositoryFactory

	// in memory
	statCacheRep persistence.MediaUserWatchStatCacheRepository
}

func NewMediaUserWatchStat(
	logger *slog.Logger,

	// repositories
	statRepo persistence.MediaUserWatchStatRepositoryFactory,

	// in memory
	statCacheRep persistence.MediaUserWatchStatCacheRepository,
) *MediaUserWatchStat {
	return &MediaUserWatchStat{
		logger:       logger,
		statRepo:     statRepo,
		statCacheRep: statCacheRep,
	}
}
