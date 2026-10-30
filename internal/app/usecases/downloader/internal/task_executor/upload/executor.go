package uploadexecutor

import (
	"log/slog"

	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	downloadstatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download_status"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task_executor/mappers"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Executor struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// Storages
	storage pstorage.DownloadsStorage

	// Services
	ffmpeg pservices.FFMpeg

	// Usecases
	download     *mediadownload.MediaDownload
	uploadStatus downloadstatus.StatusUpdater
	thumbnail    thumbnail.Thumbnail

	// Functions
	updateSystemDiskInfo func()

	// Broadcaster
	broadcaster Broadcaster
}

func NewExecutor(
	logger *slog.Logger,

	// Storages
	downloadsStorage pstorage.DownloadsStorage,

	// Services
	ffmpeg pservices.FFMpeg,

	// Usecases
	download *mediadownload.MediaDownload,
	uploadStatus downloadstatus.StatusUpdater,
	thumbnail thumbnail.Thumbnail,

	// Functions
	updateSystemDiskInfo func(),

	// Broadcaster
	broadcaster Broadcaster,
) *Executor {
	return &Executor{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// Storages
		storage: downloadsStorage,

		// Services
		ffmpeg: ffmpeg,

		// Usecases
		download:     download,
		uploadStatus: uploadStatus,
		thumbnail:    thumbnail,

		// Functions
		updateSystemDiskInfo: updateSystemDiskInfo,

		// Broadcaster
		broadcaster: broadcaster,
	}
}
