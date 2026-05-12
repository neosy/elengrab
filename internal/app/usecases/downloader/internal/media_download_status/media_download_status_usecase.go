package downloadstatus

import (
	"log/slog"

	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task_status"
	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	statussetter "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download_status/status_setter"
)

type MediaDownloadStatus struct {
	logger *slog.Logger

	// internal
	statusSetter *statussetter.MediaDownloadStatusSetter

	// usecases
	download     *mediadownload.MediaDownload
	dlTask       *dltask.DownloadTask
	dlTaskStatus *dltasktatus.DownloadTaskStatus
}

func NewMediaDownloadStatus(
	logger *slog.Logger,

	// usecases
	download *mediadownload.MediaDownload,
	dlTask *dltask.DownloadTask,
	dlTaskStatus *dltasktatus.DownloadTaskStatus,
) *MediaDownloadStatus {
	return &MediaDownloadStatus{
		logger: logger,

		// internal
		statusSetter: statussetter.NewMediaDownloadStatusSetter(logger),

		// usecases
		download:     download,
		dlTask:       dlTask,
		dlTaskStatus: dlTaskStatus,
	}
}
