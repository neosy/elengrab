package downloader

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	dlstate "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_state_cache"
	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task_status"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/file"
	filestatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/file_status"
	siteicon "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/site_icon"
	iconfetcher "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/site_icon_fetcher"
	ytchannel "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/youtube_channel"
	"github.com/neosy/elengrab/internal/app/usecases/mappers"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	"github.com/neosy/elengrab/pkg/nworkerpool"
)

type YouTubeDownloader struct {
	appCtx  context.Context
	logger  *slog.Logger
	mappers *mappers.Mappers

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

	//broadcasters
	broadcaster *broadcaster.Broadcaster

	// services
	downloaderSrv pservices.Downloader

	// Options
	downloadsDir          string
	historyMode           dtypes.HistoryMode
	deleteDuplicatesScope dtypes.UniquenessScope
	logoUpdateInterval    time.Duration
	channelUpdateInterval time.Duration
}

func NewYouTubeDownloader(
	ctx context.Context,
	logger *slog.Logger,

	// repositories
	fileRep persistence.FileRepository,
	dlTaskRep persistence.DownloadTaskRepository,
	ytChannelRep persistence.YoutubeChannelRepository,
	siteLogoRep persistence.SiteLogoRepository,

	// in memory
	downloadStateCacheRep persistence.DownloadStateCacheRepository,
	ytChannelCacheRep persistence.YoutubeChannelCacheRepository,
	siteLogoCacheRep persistence.SiteLogoCacheRepository,

	// dispetchers
	dlDispetcher nworkerpool.JobDispatcher,

	// services
	downloaderSrv pservices.Downloader,

	// options
	downloadsDir string,
	historyMode dtypes.HistoryMode,
	deleteDuplicatesScope dtypes.UniquenessScope,
	logoUpdateInterval time.Duration,
	channelUpdateInterval time.Duration,

) *YouTubeDownloader {
	dlStateCache := dlstate.NewDownloadStateCache(logger, downloadStateCacheRep)

	dlTask := dltask.NewDownloadTask(logger, dlTaskRep, dlStateCache)
	file := fileuc.NewFile(logger, fileRep, dlTask, dlStateCache)
	dlTaskStatus := dltasktatus.NewDownloadTaskStatus(logger, dlTask)

	return &YouTubeDownloader{
		appCtx:  ctx,
		logger:  logger,
		mappers: mappers.NewMappers(),

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

		//broadcasters
		broadcaster: broadcaster.NewBroadcaster(),

		// services
		downloaderSrv: downloaderSrv,

		// Options
		downloadsDir:          downloadsDir,
		historyMode:           historyMode,
		deleteDuplicatesScope: deleteDuplicatesScope,
		logoUpdateInterval:    logoUpdateInterval,
		channelUpdateInterval: channelUpdateInterval,
	}
}

func (uc *YouTubeDownloader) HistoryMode() dtypes.HistoryMode {
	return uc.historyMode
}
