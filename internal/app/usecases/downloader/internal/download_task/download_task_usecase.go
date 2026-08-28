package dltask

import (
	"log/slog"

	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadTask struct {
	logger *slog.Logger

	// repositories
	TaskRepo persistence.DownloadTaskRepositoryFactory

	// internal
	dlStateCache *dlstate.DownloadStateCache
}

func NewDownloadTask(
	logger *slog.Logger,
	taskRepo persistence.DownloadTaskRepositoryFactory,
	dlStateCache *dlstate.DownloadStateCache,
) *DownloadTask {
	return &DownloadTask{
		logger:       logger,
		TaskRepo:      taskRepo,
		dlStateCache: dlStateCache,
	}
}
