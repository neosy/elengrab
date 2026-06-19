package thumbnail

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/mappers"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail/internal/repository"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Thumbnail struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// Caches
	thumbnailFileCache persistence.ThumbnailFileCacheRepository

	// Storages
	storage pstorage.ThumbnailsStorage

	// Internal usecases
	repo *repository.ThumbnailRepository
}

// NewThumbnail is the constructor for Thumbnail use case.
func NewThumbnail(
	logger *slog.Logger,
	// Repositories
	thumbnailRep persistence.ThumbnailRepository,
	// Caches
	thumbnailCache persistence.ThumbnailCacheRepository,
	thumbnailFileCache persistence.ThumbnailFileCacheRepository,
	// Usecaes
	storage pstorage.ThumbnailsStorage,
) *Thumbnail {
	return &Thumbnail{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// Cahces
		thumbnailFileCache: thumbnailFileCache,

		// Storages
		storage: storage,

		// Internal usecases
		repo: repository.NewThumbnailRepository(logger, thumbnailRep, thumbnailCache),
	}
}
