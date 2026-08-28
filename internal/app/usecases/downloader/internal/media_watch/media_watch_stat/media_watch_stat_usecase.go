package watchstat

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaWatchStat struct {
	logger *slog.Logger

	// repositories
	statRepo persistence.MediaWatchStatRepositoryFactory

	// in memory
	statCacheRep persistence.MediaWatchStatCacheRepository
}

func NewMediaWatchStat(
	logger *slog.Logger,

	// repositories
	statRepo persistence.MediaWatchStatRepositoryFactory,

	// in memory
	statCacheRep persistence.MediaWatchStatCacheRepository,
) *MediaWatchStat {
	return &MediaWatchStat{
		logger:       logger,
		statRepo:     statRepo,
		statCacheRep: statCacheRep,
	}
}
