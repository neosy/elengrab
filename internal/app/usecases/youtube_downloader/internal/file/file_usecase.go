package fileuc

import (
	"log/slog"

	dlstate "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_state_cache"
	dltask "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type File struct {
	logger *slog.Logger

	// repositories
	fileRep persistence.FileRepository

	// internal
	dlTask       *dltask.DownloadTask
	dlStateCache *dlstate.DownloadStateCache
}

func NewFile(
	logger *slog.Logger,

	// repositories
	fileRep persistence.FileRepository,

	// usecases
	dlTask *dltask.DownloadTask,
	dlStateCache *dlstate.DownloadStateCache,
) *File {
	return &File{
		logger:       logger,
		fileRep:      fileRep,
		dlTask:       dlTask,
		dlStateCache: dlStateCache,
	}
}
