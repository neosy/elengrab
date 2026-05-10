package downloader

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	dlmigration "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/dowload_data_migration"
	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task_status"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/file"
	filestatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/file_status"
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
	dlDispetcher nworkerpool.JobDispatcher

	// internal
	file            *fileuc.File
	dlTask          *dltask.DownloadTask
	fileStatus      *filestatus.FileStatus
	dlTaskStatus    *dltasktatus.DownloadTaskStatus
	ytChannel       *ytchannel.YoutubeChannel
	siteIcon        *siteicon.SiteIcon
	dlStateCache    *dlstate.DownloadStateCache
	siteIconFetcher *iconfetcher.SiteIconFetcher
	authz           *authz.Authorization
	dlDataMigration *dlmigration.DataMigration

	// broadcasters
	broadcaster *broadcaster.Broadcaster

	// usecases
	thumbnail *thumbnail.Thumbnail

	// services
	downloaderSrv pservices.Downloader
	ffmpegSrv     pservices.FFMpeg

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
	fileRep persistence.FileRepository,
	dlTaskRep persistence.DownloadTaskRepository,
	dlDataMigration persistence.DownloadDataMigrationRepository,
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
	dlDispetcher nworkerpool.JobDispatcher,

	// usecases
	thumbnail *thumbnail.Thumbnail,

	// services
	downloaderSrv pservices.Downloader,
	ffmpegSrv pservices.FFMpeg,

	// options
	demoMode bool,
	appMode dtypes.AppMode,
	deleteDuplicatesScope dtypes.UniquenessScope,
	logoUpdateInterval time.Duration,
	channelUpdateInterval time.Duration,

) *Downloader {
	dlStateCache := dlstate.NewDownloadStateCache(logger, downloadStateCacheRep)

	dlTask := dltask.NewDownloadTask(logger, dlTaskRep, dlStateCache)
	file := fileuc.NewFile(logger, fileRep, dlTask, dlStateCache)
	dlTaskStatus := dltasktatus.NewDownloadTaskStatus(logger, dlTask)

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
		dlDispetcher: dlDispetcher,

		// internal
		file:            file,
		dlTask:          dlTask,
		fileStatus:      filestatus.NewFileStatus(logger, file, dlTask, dlTaskStatus),
		dlTaskStatus:    dlTaskStatus,
		ytChannel:       ytchannel.NewYoutubeChannel(logger, ytChannelRep, ytChannelCacheRep),
		siteIcon:        siteicon.NewSiteIcon(logger, siteLogoRep, siteLogoCacheRep),
		siteIconFetcher: iconfetcher.NewSiteIconFetcher(logger),
		authz:           authz.NewAuthorization(logger, appMode),
		dlDataMigration: dlmigration.NewdataMigration(logger, dlDataMigration),

		// broadcasters
		broadcaster: broadcaster.NewBroadcaster(),

		// usecases
		thumbnail: thumbnail,

		// services
		downloaderSrv: downloaderSrv,
		ffmpegSrv:     ffmpegSrv,

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
