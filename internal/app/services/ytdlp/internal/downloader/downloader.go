package downloader

import (
	"log/slog"
	"path"
	"path/filepath"

	ffmpegsrv "github.com/neosy/elengrab/internal/app/services/ffmpeg"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/executor"
	formatcache "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/format_cache"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/mappers"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Downloader struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	formatCache *formatcache.FormatCache
	executor    *executor.Executor

	// Storages
	storage pstorage.DownloadsStorage // Storage where downloaded files are saved

	// Services
	ffmpeg *ffmpegsrv.FFmpegService

	// Parameters
	ytDlpName      string
	ytDlpPath      string
	serviceOptions idto.ServiceOptions
}

func NewDownloader(
	logger *slog.Logger,
	ytDlpPath string,
	storage pstorage.DownloadsStorage,
	ffmpeg *ffmpegsrv.FFmpegService,
	serviceOptions idto.ServiceOptions,
) *Downloader {
	formatCache := formatcache.NewFormatCache(path.Join(storage.BasePath(), consts.YtDlpFormatCacheDir))

	return &Downloader{
		logger:  logger,
		mappers: mappers.NewMappers(),

		formatCache: formatCache,
		executor:    executor.NewExecutor(logger, storage, formatCache, ytDlpPath, serviceOptions),

		// services
		ffmpeg: ffmpeg,

		// parameters
		ytDlpName:      filepath.Base(ytDlpPath),
		ytDlpPath:      ytDlpPath,
		storage:        storage,
		serviceOptions: serviceOptions,
	}
}
