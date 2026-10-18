package downloadstatus

import (
	"log/slog"

	dtypes "github.com/neosy/elengrab/internal/domain/types"

	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task_status"
	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	statussetter "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download_status/status_setter"
	processingtask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/processing_task"
	proctaskstatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/processing_task_status"
)

type mediaDownloadStatus struct {
	logger *slog.Logger

	sourceType dtypes.MediaDownloadSourceType

	// internal
	statusSetter *statussetter.MediaDownloadStatusSetter

	// usecases
	download             *mediadownload.MediaDownload
	dlTask               *dltask.DownloadTask
	processingTask       *processingtask.ProcessingTask
	dlTaskStatus         *dltasktatus.DownloadTaskStatus
	processingTaskStatus *proctaskstatus.ProcessingTaskStatus
}

// NewStatusUpdater returns the status updater for the given source type.
// It returns an error if the source type is not supported.
func NewStatusUpdater(
	logger *slog.Logger,

	sourceType dtypes.MediaDownloadSourceType,

	// usecases
	download *mediadownload.MediaDownload,
	dlTask *dltask.DownloadTask,
	processingTask *processingtask.ProcessingTask,
	dlTaskStatus *dltasktatus.DownloadTaskStatus,
	processingTaskStatus *proctaskstatus.ProcessingTaskStatus,
) StatusUpdater {
	return &mediaDownloadStatus{
		logger: logger,

		sourceType: sourceType,

		// internal
		statusSetter: statussetter.NewMediaDownloadStatusSetter(logger),

		// usecases
		download:             download,
		dlTask:               dlTask,
		processingTask:       processingTask,
		dlTaskStatus:         dlTaskStatus,
		processingTaskStatus: processingTaskStatus,
	}
}
