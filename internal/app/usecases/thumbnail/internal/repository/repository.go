package repository

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type ThumbnailRepository struct {
	logger *slog.Logger

	// repositories
	repo persistence.ThumbnailRepositoryFactory

	// cache in memory
	cacheRepo persistence.ThumbnailCacheRepository
}

func NewThumbnailRepository(
	logger *slog.Logger,
	repo persistence.ThumbnailRepositoryFactory,
	cacheRepo persistence.ThumbnailCacheRepository,
) *ThumbnailRepository {
	return &ThumbnailRepository{
		logger:    logger,
		repo:      repo,
		cacheRepo: cacheRepo,
	}
}
