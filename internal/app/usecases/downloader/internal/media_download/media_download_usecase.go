package mediadownload

import (
	"log/slog"

	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	mediawatch "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch"
	searchindex "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/search_index"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaDownload struct {
	logger *slog.Logger

	// repositories
	downloadRepo persistence.MediaDownloadRepositoryFactory

	// caches
	downloadCacheRep persistence.MediaDownloadCacheRepository

	// internal
	dlTask       *dltask.DownloadTask
	dlStateCache *dlstate.DownloadStateCache
	mediaWatch   *mediawatch.MediaWatch
	searchIndex  *searchindex.SearchIndex

	// usecases
	thumbnail *thumbnail.Thumbnail
}

func NewMediaDownload(
	logger *slog.Logger,

	// repositories
	downloadRepo persistence.MediaDownloadRepositoryFactory,

	// caches
	downloadCacheRep persistence.MediaDownloadCacheRepository,

	// internal
	dlTask *dltask.DownloadTask,
	dlStateCache *dlstate.DownloadStateCache,
	mediaWatch *mediawatch.MediaWatch,
	searchIndex *searchindex.SearchIndex,

	// usecases
	thumbnail *thumbnail.Thumbnail,
) *MediaDownload {
	return &MediaDownload{
		logger:           logger,
		downloadRepo:     downloadRepo,
		downloadCacheRep: downloadCacheRep,

		// internal
		dlTask:       dlTask,
		dlStateCache: dlStateCache,
		mediaWatch:   mediaWatch,
		searchIndex:  searchIndex,

		// usecases
		thumbnail: thumbnail,
	}
}
