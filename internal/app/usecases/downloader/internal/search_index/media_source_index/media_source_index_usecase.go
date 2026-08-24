package sourceindex

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaSourceIndex struct {
	logger *slog.Logger

	// repositories
	indexRep persistence.MediaSourceIndexRepository
}

func NewMediaSourceIndex(
	logger *slog.Logger,

	// repositories
	indexRep persistence.MediaSourceIndexRepository,
) *MediaSourceIndex {
	return &MediaSourceIndex{
		logger:   logger,
		indexRep: indexRep,
	}
}
