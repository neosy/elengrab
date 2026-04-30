package thumbnail

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/mappers"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail/internal/store"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/ports/storage"
)

type Thumbnail struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	store   *store.ThumbnailStore
	storage storage.ThumbnailsStorage
}

// NewThumbnail is the constructor for Thumbnail use case.
func NewThumbnail(
	logger *slog.Logger,
	thumbnailRep persistence.ThumbnailRepository,
	storage storage.ThumbnailsStorage,
) *Thumbnail {
	return &Thumbnail{
		logger:  logger,
		mappers: mappers.NewMappers(),

		store:   store.NewThumbnailStore(logger, thumbnailRep),
		storage: storage,
	}
}
