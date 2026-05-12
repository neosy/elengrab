package mediadownload

import (
	"log/slog"

	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaDownload struct {
	logger *slog.Logger

	// repositories
	downloadRep persistence.MediaDownloadRepository

	// internal
	dlTask       *dltask.DownloadTask
	dlStateCache *dlstate.DownloadStateCache
}

func NewMediaDownload(
	logger *slog.Logger,

	// repositories
	downloadRep persistence.MediaDownloadRepository,

	// usecases
	dlTask *dltask.DownloadTask,
	dlStateCache *dlstate.DownloadStateCache,
) *MediaDownload {
	return &MediaDownload{
		logger:       logger,
		downloadRep:  downloadRep,
		dlTask:       dlTask,
		dlStateCache: dlStateCache,
	}
}
