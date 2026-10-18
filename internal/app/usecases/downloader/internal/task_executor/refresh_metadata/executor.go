package refreshmetaexecutor

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/channel"
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
	downloader pservices.Downloader
	ffmpeg     pservices.FFMpeg

	// Usecases
	download      *mediadownload.MediaDownload
	statusUpdater downloadstatus.StatusUpdater
	channel       *channel.Channel
	thumbnail     thumbnail.Thumbnail

	// Broadcaster
	broadcaster Broadcaster
}

func NewExecutor(
	logger *slog.Logger,

	// Storages
	downloadsStorage pstorage.DownloadsStorage,

	// Services
	downloaderSrv pservices.Downloader,
	ffmpegSrv pservices.FFMpeg,

	// Usecases
	download *mediadownload.MediaDownload,
	statusUpdater downloadstatus.StatusUpdater,
	channel *channel.Channel,
	thumbnail thumbnail.Thumbnail,

	// Broadcaster
	broadcaster Broadcaster,
) *Executor {
	return &Executor{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// Storages
		storage: downloadsStorage,

		// Services
		downloader: downloaderSrv,
		ffmpeg:     ffmpegSrv,

		// Usecases
		download:      download,
		statusUpdater: statusUpdater,
		channel:       channel,
		thumbnail:     thumbnail,

		// Broadcaster
		broadcaster: broadcaster,
	}
}
