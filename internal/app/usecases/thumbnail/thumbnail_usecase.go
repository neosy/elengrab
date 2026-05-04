package thumbnail

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/mappers"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail/internal/store"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Thumbnail struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	store   *store.ThumbnailStore
	storage pstorage.ThumbnailsStorage
}

// NewThumbnail is the constructor for Thumbnail use case.
func NewThumbnail(
	logger *slog.Logger,
	thumbnailRep persistence.ThumbnailRepository,
	thumbnailCacheRep persistence.ThumbnailCacheRepository,
	storage pstorage.ThumbnailsStorage,
) *Thumbnail {
	return &Thumbnail{
		logger:  logger,
		mappers: mappers.NewMappers(),

		store:   store.NewThumbnailStore(logger, thumbnailRep, thumbnailCacheRep),
		storage: storage,
	}
}
