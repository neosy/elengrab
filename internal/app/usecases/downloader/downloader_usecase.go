package downloader

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	dlexecutor "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_executor"
	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task_status"
	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	downloadstatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download_status"
	mediawatch "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch"
	searchindex "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/search_index"
	siteicon "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/site_icon"
	ytchannel "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/youtube_channel"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/app/usecases/mappers"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type downloader struct {
	appCtx  context.Context
	logger  *slog.Logger
	mappers *mappers.Mappers

	// Storages
	downloadsStorage pstorage.DownloadsStorage

	systemInfoStore systemInfoStore

	// dispetchers
	downloadDispatcher  workerpool.JobDispatcher
	operationDispatcher workerpool.JobDispatcher

	// internal
	download       *mediadownload.MediaDownload
	dlTask         *dltask.DownloadTask
	downloadStatus *downloadstatus.MediaDownloadStatus
	dlTaskStatus   *dltasktatus.DownloadTaskStatus
	ytChannel      *ytchannel.YoutubeChannel
	siteIcon       *siteicon.SiteIcon
	authz          *authz.Authorization
	mediaWatch     *mediawatch.MediaWatch
	searchIndex    *searchindex.SearchIndex
	dlExecutor     *dlexecutor.Executor

	// broadcasters
	broadcaster *broadcaster.Broadcaster

	// usecases
	thumbnail *thumbnail.Thumbnail

	// services
	downloaderSrv pservices.Downloader
	ffmpegSrv     pservices.FFMpeg
	authSrv       pservices.AuthService

	// Options
	demoMode                        bool
	appMode                         dtypes.AppMode
	deleteDuplicatesUniquenessScope dtypes.UniquenessScope
}

func NewDownloader(
	ctx context.Context,
	logger *slog.Logger,

	// repositories
	downloadRepo persistence.MediaDownloadRepositoryFactory,
	dlTaskRepo persistence.DownloadTaskRepositoryFactory,
	watchEventRepo persistence.MediaWatchEventRepositoryFactory,
	userWatchChunkRepo persistence.MediaUserWatchChunkRepositoryFactory,
	userWatchStatRepo persistence.MediaUserWatchStatRepositoryFactory,
	watchStatRepo persistence.MediaWatchStatRepositoryFactory,
	userWatchPositionRepo persistence.MediaUserWatchPositionRepositoryFactory,
	sourceIndexRepo persistence.MediaSourceIndexRepositoryFactory,
	ytChannelRepo persistence.YoutubeChannelRepositoryFactory,
	siteLogoRepo persistence.SiteLogoRepositoryFactory,

	// in memory
	mediaDownloadCacheRep persistence.MediaDownloadCacheRepository,
	downloadStateCacheRep persistence.DownloadStateCacheRepository,
	mediaUserWatchStatCacheRep persistence.MediaUserWatchStatCacheRepository,
	mediaWatchStatCacheRep persistence.MediaWatchStatCacheRepository,
	mediaUserWatchPositionCacheRep persistence.MediaUserWatchPositionCacheRepository,
	ytChannelCacheRep persistence.YoutubeChannelCacheRepository,
	siteLogoCacheRep persistence.SiteLogoCacheRepository,

	// storages
	downloadsStorage pstorage.DownloadsStorage,

	// dispetchers
	downloadDispatcher workerpool.JobDispatcher,
	operationDispatcher workerpool.JobDispatcher,
	watchEventDispatcher workerpool.JobDispatcher,

	// usecases
	thumbnail *thumbnail.Thumbnail,

	// services
	downloaderSrv pservices.Downloader,
	ffmpegSrv pservices.FFMpeg,
	authSrv pservices.AuthService,

	// options
	demoMode bool,
	appMode dtypes.AppMode,
	deleteDuplicatesUniquenessScope dtypes.UniquenessScope,
	logoUpdateInterval time.Duration,
	channelUpdateInterval time.Duration,
) *downloader {
	dlStateCache := dlstate.NewDownloadStateCache(logger, downloadStateCacheRep)
	dlTask := dltask.NewDownloadTask(logger, dlTaskRepo, dlStateCache)
	dlTaskStatus := dltasktatus.NewDownloadTaskStatus(logger, dlTask)

	authz := authz.NewAuthorization(logger, appMode)

	siteIcon := siteicon.NewSiteIcon(logger, siteLogoRepo, siteLogoCacheRep)
	ytChannel := ytchannel.NewYoutubeChannel(logger, ytChannelRepo, ytChannelCacheRep)

	searchIndex := searchindex.NewSearchIndex(logger, sourceIndexRepo)

	downloader := &downloader{
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

		// Dispetchers
		downloadDispatcher:  downloadDispatcher,
		operationDispatcher: operationDispatcher,

		// Internal
		dlTask:       dlTask,
		dlTaskStatus: dlTaskStatus,
		ytChannel:    ytChannel,
		siteIcon:     siteIcon,
		authz:        authz,

		// broadcasters
		broadcaster: broadcaster.NewBroadcaster(authz),

		// usecases
		thumbnail:   thumbnail,
		searchIndex: searchIndex,

		// services
		downloaderSrv: downloaderSrv,
		ffmpegSrv:     ffmpegSrv,
		authSrv:       authSrv,

		// Options
		demoMode:                        demoMode,
		appMode:                         appMode,
		deleteDuplicatesUniquenessScope: deleteDuplicatesUniquenessScope,
	}

	mediaWatch := mediawatch.NewMediaWatch(
		logger,

		watchEventRepo, userWatchChunkRepo, userWatchStatRepo, watchStatRepo, userWatchPositionRepo,

		mediaUserWatchStatCacheRep,
		mediaWatchStatCacheRep,
		mediaUserWatchPositionCacheRep,

		watchEventDispatcher,

		downloader.broadcastWatchStatsUpdatedToAuth,
		downloader.broadcastWatchPositionUpdatedToAuth,
	)

	download := mediadownload.NewMediaDownload(
		logger,

		downloadRepo,
		mediaDownloadCacheRep,

		dlTask,
		dlStateCache,
		mediaWatch,
		searchIndex,

		thumbnail,
	)

	downloadStatus := downloadstatus.NewMediaDownloadStatus(logger, download, dlTask, dlTaskStatus)

	dlExecutor := dlexecutor.NewExecutor(
		ctx,
		logger,

		// Storages
		downloadsStorage,

		// Services
		downloaderSrv,
		ffmpegSrv,

		// Usecases
		download,
		downloadStatus,
		siteIcon,
		ytChannel,
		thumbnail,

		// Broadcaster
		dlexecutor.Broadcaster{
			DownloadUpdate:         downloader.broadcastDownloadUpdate,
			DownloadProgressUpdate: downloader.broadcastDownloadProgressUpdate,
		},

		// Options
		logoUpdateInterval,
		channelUpdateInterval,
	)

	// Internal
	downloader.download = download
	downloader.downloadStatus = downloadStatus
	downloader.mediaWatch = mediaWatch
	downloader.dlExecutor = dlExecutor

	return downloader
}

func (uc *downloader) AppMode() dtypes.AppMode {
	return uc.appMode
}

func (uc *downloader) DemoMode() bool {
	return uc.demoMode
}
