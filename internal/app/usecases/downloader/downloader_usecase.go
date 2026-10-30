package downloader

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	uchannel "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/channel"
	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task_status"
	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	downloadstatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download_status"
	mediawatch "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch"
	processingtask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/processing_task"
	proctaskstatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/processing_task_status"
	searchindex "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/search_index"
	siteicon "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/site_icon"
	taskexecutor "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task_executor"
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
	downloadDispatcher        workerpool.JobDispatcher
	mediaProcessingDispatcher workerpool.JobDispatcher
	operationDispatcher       workerpool.JobDispatcher

	// internal
	download             *mediadownload.MediaDownload
	dlTask               *dltask.DownloadTask
	processingTask       *processingtask.ProcessingTask
	downloadStatus       downloadstatus.StatusUpdater
	uploadStatus         downloadstatus.StatusUpdater
	dlTaskStatus         *dltasktatus.DownloadTaskStatus
	processingTaskStatus *proctaskstatus.ProcessingTaskStatus
	channel              *uchannel.Channel
	siteIcon             *siteicon.SiteIcon
	authz                *authz.Authorization
	mediaWatch           *mediawatch.MediaWatch
	searchIndex          *searchindex.SearchIndex
	taskExecutors        *taskexecutor.TaskExecutors

	// broadcasters
	broadcaster *broadcaster.Broadcaster

	// usecases
	thumbnail thumbnail.Thumbnail

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
	processingTaskRepo persistence.MediaProcessingTaskRepositoryFactory,
	watchEventRepo persistence.MediaWatchEventRepositoryFactory,
	userWatchChunkRepo persistence.MediaUserWatchChunkRepositoryFactory,
	userWatchStatRepo persistence.MediaUserWatchStatRepositoryFactory,
	watchStatRepo persistence.MediaWatchStatRepositoryFactory,
	userWatchPositionRepo persistence.MediaUserWatchPositionRepositoryFactory,
	sourceIndexRepo persistence.MediaSourceIndexRepositoryFactory,
	ytChannelRepo persistence.ChannelRepositoryFactory,
	siteLogoRepo persistence.SiteLogoRepositoryFactory,

	// in memory
	mediaDownloadCacheRep persistence.MediaDownloadCacheRepository,
	downloadStateCacheRep persistence.DownloadStateCacheRepository,
	mediaUserWatchStatCacheRep persistence.MediaUserWatchStatCacheRepository,
	mediaWatchStatCacheRep persistence.MediaWatchStatCacheRepository,
	mediaUserWatchPositionCacheRep persistence.MediaUserWatchPositionCacheRepository,
	ytChannelCacheRep persistence.ChannelCacheRepository,
	siteLogoCacheRep persistence.SiteLogoCacheRepository,

	// storages
	downloadsStorage pstorage.DownloadsStorage,

	// dispetchers
	downloadDispatcher workerpool.JobDispatcher,
	mediaProcessingDispatcher workerpool.JobDispatcher,
	operationDispatcher workerpool.JobDispatcher,
	watchEventDispatcher workerpool.JobDispatcher,

	// usecases
	thumbnail thumbnail.Thumbnail,

	// services
	downloaderSrv pservices.Downloader,
	ffmpegSrv pservices.FFMpeg,
	authSrv pservices.AuthService,

	// options
	demoMode bool,
	appMode dtypes.AppMode,
	deleteDuplicatesUniquenessScope dtypes.UniquenessScope,
	logoUpdateInterval time.Duration,
) *downloader {
	dlStateCache := dlstate.NewDownloadStateCache(logger, downloadStateCacheRep)
	dlTask := dltask.NewDownloadTask(logger, dlTaskRepo, dlStateCache)
	processingTask := processingtask.NewProcessingTask(logger, processingTaskRepo)
	dlTaskStatus := dltasktatus.NewDownloadTaskStatus(logger, dlTask)
	processingTaskStatus := proctaskstatus.NewProcessingTaskStatus(logger, processingTask)

	authz := authz.NewAuthorization(logger, appMode)

	siteIcon := siteicon.NewSiteIcon(logger, siteLogoRepo, siteLogoCacheRep)
	channel := uchannel.NewChannel(logger, ytChannelRepo, ytChannelCacheRep)

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
		downloadDispatcher:        downloadDispatcher,
		mediaProcessingDispatcher: mediaProcessingDispatcher,
		operationDispatcher:       operationDispatcher,

		// Internal
		dlTask:               dlTask,
		processingTask:       processingTask,
		dlTaskStatus:         dlTaskStatus,
		processingTaskStatus: processingTaskStatus,
		channel:              channel,
		siteIcon:             siteIcon,
		authz:                authz,

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

		downloader.onWatchUserStatsUpdated,
		downloader.onWatchStatsUpdated,
		downloader.onWatchUserPositionUpdated,
	)

	download := mediadownload.NewMediaDownload(
		logger,

		downloadRepo,
		mediaDownloadCacheRep,

		dlTask,
		dlStateCache,
		processingTask,
		mediaWatch,
		searchIndex,

		thumbnail,
	)

	downloadStatus := downloadstatus.NewStatusUpdater(
		logger,
		dtypes.MediaDownloadSourceTypeURL,
		download,
		dlTask, processingTask,
		dlTaskStatus, processingTaskStatus,
	)

	uploadStatus := downloadstatus.NewStatusUpdater(
		logger,
		dtypes.MediaDownloadSourceTypeUpload,
		download,
		dlTask, processingTask,
		dlTaskStatus, processingTaskStatus,
	)

	taskExecutors := taskexecutor.NewTaskExecutors(
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
		uploadStatus,
		siteIcon,
		channel,
		thumbnail,

		// Functions
		downloader.UpdateSystemInfo,

		// Broadcaster
		taskexecutor.Broadcaster{
			MediaUpdated:            downloader.broadcastMediaDownloadUpdated,
			DownloadProgressUpdated: downloader.broadcastDownloadProgressUpdated,
			StartRefreshing:         downloader.broadcastMediaDownloadStartRefreshing,
		},

		// Options
		logoUpdateInterval,
	)

	// Internal
	downloader.download = download
	downloader.downloadStatus = downloadStatus
	downloader.uploadStatus = uploadStatus
	downloader.mediaWatch = mediaWatch
	downloader.taskExecutors = taskExecutors

	return downloader
}

func (uc *downloader) AppMode() dtypes.AppMode {
	return uc.appMode
}

func (uc *downloader) DemoMode() bool {
	return uc.demoMode
}
