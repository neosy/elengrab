package uwatchstat

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaUserWatchStat struct {
	logger *slog.Logger

	// repositories
	statRep persistence.MediaUserWatchStatRepository
}

func NewMediaUserWatchStat(
	logger *slog.Logger,

	// repositories
	statRep persistence.MediaUserWatchStatRepository,
) *MediaUserWatchStat {
	return &MediaUserWatchStat{
		logger:  logger,
		statRep: statRep,
	}
}
