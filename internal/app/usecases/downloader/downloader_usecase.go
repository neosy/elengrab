package downloader

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	dlmigration "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_data_migration"
	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task_status"
	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	downloadstatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download_status"
	siteicon "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/site_icon"
	iconfetcher "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/site_icon_fetcher"
	ytchannel "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/youtube_channel"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/app/usecases/mappers"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Downloader struct {
	appCtx  context.Context
	logger  *slog.Logger
	mappers *mappers.Mappers

	// Storages
	downloadsStorage pstorage.DownloadsStorage

	systemInfoStore systemInfoStore

	// dispetchers
	downloadDispatcher  nworkerpool.JobDispatcher
	operationDispatcher nworkerpool.JobDispatcher

	// internal
	download          *mediadownload.MediaDownload
	dlTask            *dltask.DownloadTask
	downloadStatus    *downloadstatus.MediaDownloadStatus
	dlTaskStatus      *dltasktatus.DownloadTaskStatus
	ytChannel         *ytchannel.YoutubeChannel
	siteIcon          *siteicon.SiteIcon
	dlStateCache      *dlstate.DownloadStateCache
	siteIconFetcher   *iconfetcher.SiteIconFetcher
	authz             *authz.Authorization
	downloadMigration *dlmigration.DownloadMigration

	// broadcasters
	broadcaster *broadcaster.Broadcaster

	// usecases
	thumbnail *thumbnail.Thumbnail

	// services
	downloaderSrv pservices.Downloader
	ffmpegSrv     pservices.FFMpeg
	authSrv       pservices.AuthService

	// Options
	demoMode              bool
	appMode               dtypes.AppMode
	deleteDuplicatesScope dtypes.UniquenessScope
	logoUpdateInterval    time.Duration
	channelUpdateInterval time.Duration
}

func NewDownloader(
	ctx context.Context,
	logger *slog.Logger,

	// repositories
	downloadRep persistence.MediaDownloadRepository,
	dlTaskRep persistence.DownloadTaskRepository,
	downloadMigrationRep persistence.DownloadDataMigrationRepository,
	ytChannelRep persistence.YoutubeChannelRepository,
	siteLogoRep persistence.SiteLogoRepository,
	thumbnailRep persistence.ThumbnailRepository,

	// in memory
	downloadStateCacheRep persistence.DownloadStateCacheRepository,
	ytChannelCacheRep persistence.YoutubeChannelCacheRepository,
	siteLogoCacheRep persistence.SiteLogoCacheRepository,

	// storages
	downloadsStorage pstorage.DownloadsStorage,

	// dispetchers
	downloadDispatcher nworkerpool.JobDispatcher,
	operationDispatcher nworkerpool.JobDispatcher,

	// usecases
	thumbnail *thumbnail.Thumbnail,

	// services
	downloaderSrv pservices.Downloader,
	ffmpegSrv pservices.FFMpeg,
	authSrv pservices.AuthService,

	// options
	demoMode bool,
	appMode dtypes.AppMode,
	deleteDuplicatesScope dtypes.UniquenessScope,
	logoUpdateInterval time.Duration,
	channelUpdateInterval time.Duration,

) *Downloader {
	dlStateCache := dlstate.NewDownloadStateCache(logger, downloadStateCacheRep)

	dlTask := dltask.NewDownloadTask(logger, dlTaskRep, dlStateCache)
	download := mediadownload.NewMediaDownload(logger, downloadRep, dlTask, dlStateCache)
	dlTaskStatus := dltasktatus.NewDownloadTaskStatus(logger, dlTask)

	authz := authz.NewAuthorization(logger, appMode)

	return &Downloader{
		appCtx:  ctx,
		logger:  logger,
		mappers: mappers.NewMappers(),

		// Storages
		downloadsStorage: downloadsStorage,

		systemInfoStore: systemInfoStore{
			data: dto.SystemInfoResponse{
				AppVersion: iconfig.AppVersion,
			},
		},

		// Cache
		dlStateCache: dlStateCache,

		// dispetchers
		downloadDispatcher:  downloadDispatcher,
		operationDispatcher: operationDispatcher,

		// internal
		download:          download,
		dlTask:            dlTask,
		downloadStatus:    downloadstatus.NewMediaDownloadStatus(logger, download, dlTask, dlTaskStatus),
		dlTaskStatus:      dlTaskStatus,
		ytChannel:         ytchannel.NewYoutubeChannel(logger, ytChannelRep, ytChannelCacheRep),
		siteIcon:          siteicon.NewSiteIcon(logger, siteLogoRep, siteLogoCacheRep),
		siteIconFetcher:   iconfetcher.NewSiteIconFetcher(logger),
		authz:             authz,
		downloadMigration: dlmigration.NewDownloadMigration(logger, downloadMigrationRep),

		// broadcasters
		broadcaster: broadcaster.NewBroadcaster(authz),

		// usecases
		thumbnail: thumbnail,

		// services
		downloaderSrv: downloaderSrv,
		ffmpegSrv:     ffmpegSrv,
		authSrv:       authSrv,

		// Options
		demoMode:              demoMode,
		appMode:               appMode,
		deleteDuplicatesScope: deleteDuplicatesScope,
		logoUpdateInterval:    logoUpdateInterval,
		channelUpdateInterval: channelUpdateInterval,
	}
}

func (uc *Downloader) AppMode() dtypes.AppMode {
	return uc.appMode
}

func (uc *Downloader) DemoMode() bool {
	return uc.demoMode
}
