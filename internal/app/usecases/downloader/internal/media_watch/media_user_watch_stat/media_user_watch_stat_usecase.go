package uwatchstat

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaUserWatchStat struct {
	logger *slog.Logger

	// repositories
	statRep persistence.MediaUserWatchStatRepository

	// in memory
	statCacheRep persistence.MediaUserWatchStatCacheRepository
}

func NewMediaUserWatchStat(
	logger *slog.Logger,

	// repositories
	statRep persistence.MediaUserWatchStatRepository,

	// in memory
	statCacheRep persistence.MediaUserWatchStatCacheRepository,
) *MediaUserWatchStat {
	return &MediaUserWatchStat{
		logger:       logger,
		statRep:      statRep,
		statCacheRep: statCacheRep,
	}
}
