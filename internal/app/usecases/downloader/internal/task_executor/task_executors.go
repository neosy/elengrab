package taskexecutor

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/channel"
	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	downloadstatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download_status"
	siteicon "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/site_icon"
	dlexecutor "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task_executor/download"
	refreshmetaexecutor "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task_executor/refresh_metadata"
	uploadexecutor "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task_executor/upload"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type TaskExecutors struct {
	Download        *dlexecutor.Executor
	Upload          *uploadexecutor.Executor
	RefreshMetadata *refreshmetaexecutor.Executor
}

func NewTaskExecutors(
	ctx context.Context,
	logger *slog.Logger,

	// Storages
	downloadsStorage pstorage.DownloadsStorage,

	// Services
	downloader pservices.Downloader,
	ffmpeg pservices.FFMpeg,

	// Usecases
	download *mediadownload.MediaDownload,
	downloadStatus downloadstatus.StatusUpdater,
	uploadStatus downloadstatus.StatusUpdater,
	siteIcon *siteicon.SiteIcon,
	channel *channel.Channel,
	thumbnail thumbnail.Thumbnail,

	// Functions
	updateSystemDiskInfo func(),

	// Broadcaster
	broadcaster Broadcaster,

	// Options
	logoUpdateInterval time.Duration,
) *TaskExecutors {
	dlExecutor := dlexecutor.NewExecutor(
		ctx,
		logger,

		// Storages
		downloadsStorage,

		// Services
		downloader,
		ffmpeg,

		// Usecases
		download,
		downloadStatus,
		siteIcon,
		channel,
		thumbnail,

		// Functions
		updateSystemDiskInfo,

		// Broadcaster
		dlexecutor.Broadcaster{
			MediaUpdated:            broadcaster.MediaUpdated,
			DownloadProgressUpdated: broadcaster.DownloadProgressUpdated,
		},

		// Options
		logoUpdateInterval,
	)

	uploadExecutor := uploadexecutor.NewExecutor(
		logger,

		// Storages,
		downloadsStorage,

		// Services,
		ffmpeg,

		// Usecases
		download,
		uploadStatus,
		thumbnail,

		// Functions
		updateSystemDiskInfo,

		// Broadcaster
		uploadexecutor.Broadcaster{
			MediaUpdated:    broadcaster.MediaUpdated,
			StartRefreshing: broadcaster.StartRefreshing,
		},
	)

	refreshMetadataExecutor := refreshmetaexecutor.NewExecutor(
		logger,

		// Storages,
		downloadsStorage,

		// Services,
		downloader,
		ffmpeg,

		// Usecases
		download,
		downloadStatus,
		channel,
		thumbnail,

		// Broadcaster
		refreshmetaexecutor.Broadcaster{
			MediaUpdated: broadcaster.MediaUpdated,
		},
	)

	return &TaskExecutors{
		Download:        dlExecutor,
		Upload:          uploadExecutor,
		RefreshMetadata: refreshMetadataExecutor,
	}
}
