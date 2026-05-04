package store

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type ThumbnailStore struct {
	logger *slog.Logger

	// repositories
	thumbnailRep persistence.ThumbnailRepository

	// cache in memory
	thumbnailCacheRep persistence.ThumbnailCacheRepository
}

func NewThumbnailStore(
	logger *slog.Logger,
	thumbnailRep persistence.ThumbnailRepository,
	thumbnailCacheRep persistence.ThumbnailCacheRepository,
) *ThumbnailStore {
	return &ThumbnailStore{
		logger:            logger,
		thumbnailRep:      thumbnailRep,
		thumbnailCacheRep: thumbnailCacheRep,
	}
}
