package sourceindex

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaSourceIndex struct {
	logger *slog.Logger

	// repositories
	indexRepo persistence.MediaSourceIndexRepositoryFactory
}

func NewMediaSourceIndex(
	logger *slog.Logger,

	// repositories
	indexRepo persistence.MediaSourceIndexRepositoryFactory,
) *MediaSourceIndex {
	return &MediaSourceIndex{
		logger:    logger,
		indexRepo: indexRepo,
	}
}
