package watchchunk

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaWatchChunk struct {
	logger *slog.Logger

	// repositories
	chunkRep persistence.MediaWatchChunkRepository
}

func NewMediaWatchChunk(
	logger *slog.Logger,

	// repositories
	chunkRep persistence.MediaWatchChunkRepository,
) *MediaWatchChunk {
	return &MediaWatchChunk{
		logger:   logger,
		chunkRep: chunkRep,
	}
}
