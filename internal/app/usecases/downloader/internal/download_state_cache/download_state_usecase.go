package dlstate

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadStateCache struct {
	logger *slog.Logger

	// repositories
	stateCacheRep persistence.DownloadStateCacheRepository
}

func NewDownloadStateCache(
	logger *slog.Logger,
	stateCacheRep persistence.DownloadStateCacheRepository,
) *DownloadStateCache {
	return &DownloadStateCache{
		logger:        logger,
		stateCacheRep: stateCacheRep,
	}
}
