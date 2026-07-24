package mediadownload

import (
	"log/slog"

	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	mediawatch "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaDownload struct {
	logger *slog.Logger

	// repositories
	downloadRep persistence.MediaDownloadRepository

	// caches
	downloadCacheRep persistence.MediaDownloadCacheRepository

	// internal
	dlTask       *dltask.DownloadTask
	dlStateCache *dlstate.DownloadStateCache
	mediaWatch   *mediawatch.MediaWatch

	// usecases
	thumbnail *thumbnail.Thumbnail
}

func NewMediaDownload(
	logger *slog.Logger,

	// repositories
	downloadRep persistence.MediaDownloadRepository,

	// caches
	downloadCacheRep persistence.MediaDownloadCacheRepository,

	// internal
	dlTask *dltask.DownloadTask,
	dlStateCache *dlstate.DownloadStateCache,
	mediaWatch *mediawatch.MediaWatch,

	// usecases
	thumbnail *thumbnail.Thumbnail,
) *MediaDownload {
	return &MediaDownload{
		logger:           logger,
		downloadRep:      downloadRep,
		downloadCacheRep: downloadCacheRep,

		// internal
		dlTask:       dlTask,
		dlStateCache: dlStateCache,
		mediaWatch:   mediaWatch,

		// usecases
		thumbnail: thumbnail,
	}
}
