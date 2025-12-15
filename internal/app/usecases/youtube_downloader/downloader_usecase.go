package ytdownloader

import (
	"context"
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/mappers"
	dlstate "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_state"
	dltask "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task_status"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/file"
	filestatus "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/file_status"
	ytchannel "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/youtube_channel"
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
	file         *fileuc.File
	dlTask       *dltask.DownloadTask
	fileStatus   *filestatus.FileStatus
	dlTaskStatus *dltasktatus.DownloadTaskStatus
	ytChannel    *ytchannel.YoutubeChannel
	dlState      *dlstate.DownloadState

	// services
	downloaderSrv pservices.YouTubeDownloader

	// Options
	downloadsDir string
	loadHistory  bool
}

func NewYouTubeDownloader(
	ctx context.Context,
	logger *slog.Logger,

	// repositories
	fileRep persistence.FileRepository,
	dlTaskRep persistence.DownloadTaskRepository,
	ytChannelRep persistence.YoutubeChannelRepository,
	downloadStateRep persistence.DownloadStateRepository,

	// dispetchers
	dlDispetcher nworkerpool.JobDispatcher,

	// services
	downloaderSrv pservices.YouTubeDownloader,

	// options
	downloadsDir string,
	loadHistory bool,
) *YouTubeDownloader {
	dlTask := dltask.NewDownloadTask(logger, dlTaskRep)
	file := fileuc.NewFile(logger, fileRep, dlTask)
	dlTaskStatus := dltasktatus.NewDownloadTaskStatus(logger, dlTaskRep, dlTask)

	return &YouTubeDownloader{
		appCtx:  ctx,
		logger:  logger,
		mappers: mappers.NewMappers(),

		// in memory
		dlState: dlstate.NewDownloadState(logger, downloadStateRep),

		// dispetchers
		dlDispetcher: dlDispetcher,

		// internal
		file:         file,
		dlTask:       dlTask,
		fileStatus:   filestatus.NewFileStatus(logger, fileRep, file, dlTask, dlTaskStatus),
		dlTaskStatus: dlTaskStatus,
		ytChannel:    ytchannel.NewYoutubeChannel(logger, ytChannelRep),

		// services
		downloaderSrv: downloaderSrv,

		// Options
		downloadsDir: downloadsDir,
		loadHistory:  loadHistory,
	}
}
