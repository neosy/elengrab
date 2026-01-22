package dlstate

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadStateCache struct {
	logger *slog.Logger

	// repositories
	stateRep persistence.DownloadStateRepository
}

func NewDownloadStateCache(
	logger *slog.Logger,
	stateRep persistence.DownloadStateRepository,
) *DownloadStateCache {
	return &DownloadStateCache{
		logger:   logger,
		stateRep: stateRep,
	}
}
