package store

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type ThumbnailStore struct {
	logger *slog.Logger

	// repositories
	thumbnailRep persistence.ThumbnailRepository
}

func NewThumbnailStore(
	logger *slog.Logger,
	thumbnailRep persistence.ThumbnailRepository,
) *ThumbnailStore {
	return &ThumbnailStore{
		logger:       logger,
		thumbnailRep: thumbnailRep,
	}
}
