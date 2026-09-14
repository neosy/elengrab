package dlexecutor

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/channel"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_executor/mappers"
	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	downloadstatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download_status"
	siteicon "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/site_icon"
	iconfetcher "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/site_icon_fetcher"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Broadcaster struct {
	DownloadUpdate         func(context.Context, uuid.UUID)
	DownloadProgressUpdate func(context.Context, uuid.UUID)
}

type Executor struct {
	appCtx  context.Context
	logger  *slog.Logger
	mappers *mappers.Mappers

	// Storages
	downloadsStorage pstorage.DownloadsStorage

	// Services
	downloaderSrv pservices.Downloader
	ffmpegSrv     pservices.FFMpeg

	// Usecases
	download        *mediadownload.MediaDownload
	downloadStatus  *downloadstatus.MediaDownloadStatus
	siteIcon        *siteicon.SiteIcon
	siteIconFetcher *iconfetcher.SiteIconFetcher
	channel         *channel.Channel
	thumbnail       thumbnail.Thumbnail

	// Broadcaster
	broadcaster Broadcaster

	// Options
	logoUpdateInterval time.Duration
}

func NewExecutor(
	appCtx context.Context,
	logger *slog.Logger,

	// Storages
	downloadsStorage pstorage.DownloadsStorage,

	// Services
	downloaderSrv pservices.Downloader,
	ffmpegSrv pservices.FFMpeg,

	// Usecases
	download *mediadownload.MediaDownload,
	downloadStatus *downloadstatus.MediaDownloadStatus,
	siteIcon *siteicon.SiteIcon,
	channel *channel.Channel,
	thumbnail thumbnail.Thumbnail,

	// Broadcaster
	broadcaster Broadcaster,

	// Options
	logoUpdateInterval time.Duration,
) *Executor {
	return &Executor{
		appCtx:  appCtx,
		logger:  logger,
		mappers: mappers.NewMappers(),

		// Storages
		downloadsStorage: downloadsStorage,

		// Services
		downloaderSrv: downloaderSrv,
		ffmpegSrv:     ffmpegSrv,

		// Usecases
		download:        download,
		downloadStatus:  downloadStatus,
		siteIcon:        siteIcon,
		siteIconFetcher: iconfetcher.NewSiteIconFetcher(logger),
		channel:         channel,
		thumbnail:       thumbnail,

		// Broadcaster
		broadcaster: broadcaster,

		// Options
		logoUpdateInterval: logoUpdateInterval,
	}
}
