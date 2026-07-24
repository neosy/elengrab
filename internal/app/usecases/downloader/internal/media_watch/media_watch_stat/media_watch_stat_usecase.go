package watchstat

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaWatchStat struct {
	logger *slog.Logger

	// repositories
	statRep persistence.MediaWatchStatRepository

	// in memory
	statCacheRep persistence.MediaWatchStatCacheRepository
}

func NewMediaWatchStat(
	logger *slog.Logger,

	// repositories
	statRep persistence.MediaWatchStatRepository,

	// in memory
	statCacheRep persistence.MediaWatchStatCacheRepository,
) *MediaWatchStat {
	return &MediaWatchStat{
		logger:       logger,
		statRep:      statRep,
		statCacheRep: statCacheRep,
	}
}
