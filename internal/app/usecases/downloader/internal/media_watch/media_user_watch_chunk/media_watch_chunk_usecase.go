package uwatchchunk

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaUserWatchChunk struct {
	logger *slog.Logger

	// repositories
	chunkRep persistence.MediaUserWatchChunkRepository
}

func NewMediaUserWatchChunk(
	logger *slog.Logger,

	// repositories
	chunkRep persistence.MediaUserWatchChunkRepository,
) *MediaUserWatchChunk {
	return &MediaUserWatchChunk{
		logger:   logger,
		chunkRep: chunkRep,
	}
}
