package dltask

import (
	"log/slog"

	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadTask struct {
	logger *slog.Logger

	// repositories
	TaskRep persistence.DownloadTaskRepository

	// internal
	dlStateCache *dlstate.DownloadStateCache
}

func NewDownloadTask(
	logger *slog.Logger,
	taskRep persistence.DownloadTaskRepository,
	dlStateCache *dlstate.DownloadStateCache,
) *DownloadTask {
	return &DownloadTask{
		logger:       logger,
		TaskRep:      taskRep,
		dlStateCache: dlStateCache,
	}
}
