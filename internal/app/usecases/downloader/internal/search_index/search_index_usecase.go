package searchindex

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/search_index/mappers"
	sourceindex "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/search_index/media_source_index"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type SearchIndex struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// Internal usecase
	searchIndex *sourceindex.MediaSourceIndex
}

func NewSearchIndex(
	logger *slog.Logger,

	// Repositories
	sourceIndexRep persistence.MediaSourceIndexRepository,
) *SearchIndex {
	return &SearchIndex{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// Internal usecase
		searchIndex: sourceindex.NewMediaSourceIndex(logger, sourceIndexRep),
	}
}
