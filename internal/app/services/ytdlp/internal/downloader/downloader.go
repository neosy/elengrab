package downloader

import (
	"log/slog"
	"path"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/executor"
	formatcache "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/format_cache"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/mappers"
)

type Downloader struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	formatCache *formatcache.FormatCache
	executor    *executor.Executor

	// parameters
	ytDlpName      string
	ytDlpPath      string
	downloadsDir   string // Directory where downloaded files are saved
	serviceOptions *dto.Options
}

func NewDownloader(
	logger *slog.Logger,
	ytDlpPath string,
	downloadsDir string,
	serviceOptions *dto.Options,
) *Downloader {
	formatCache := formatcache.NewFormatCache(path.Join(downloadsDir, consts.YtDlpFormatCacheDir))

	return &Downloader{
		logger:  logger,
		mappers: mappers.NewMappers(),

		formatCache: formatCache,
		executor:    executor.NewExecutor(logger, formatCache, ytDlpPath, serviceOptions),

		// parameters
		ytDlpName:      filepath.Base(ytDlpPath),
		ytDlpPath:      ytDlpPath,
		downloadsDir:   downloadsDir,
		serviceOptions: serviceOptions,
	}
}
