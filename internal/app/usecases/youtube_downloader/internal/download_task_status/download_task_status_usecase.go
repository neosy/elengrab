package dltaskstatus

import (
	"log/slog"

	downloadtask "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task"
	statussetter "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task_status/status_setter"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadTaskStatus struct {
	logger  *slog.Logger
	taskRep persistence.DownloadTaskRepository

	// internal
	statusSetter *statussetter.DownloadTaskStatusSetter

	// usecases
	downloadTask *downloadtask.DownloadTask
}

func NewDownloadTaskStatus(
	logger *slog.Logger,

	// repositories
	taskRep persistence.DownloadTaskRepository,

	// usecases
	downloadTask *downloadtask.DownloadTask,
) *DownloadTaskStatus {
	return &DownloadTaskStatus{
		logger:  logger,
		taskRep: taskRep,

		// internal
		statusSetter: statussetter.NewFDownloadTaskStatusSetter(logger),

		// usecases
		downloadTask: downloadTask,
	}
}
