package uwatchchunk

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaUserWatchChunk struct {
	logger *slog.Logger

	// repositories
	chunkRepo persistence.MediaUserWatchChunkRepositoryFactory
}

func NewMediaUserWatchChunk(
	logger *slog.Logger,

	// repositories
	chunkRepo persistence.MediaUserWatchChunkRepositoryFactory,
) *MediaUserWatchChunk {
	return &MediaUserWatchChunk{
		logger:    logger,
		chunkRepo: chunkRepo,
	}
}
