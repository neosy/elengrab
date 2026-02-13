package core

import (
	"log/slog"
	"path"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/downloader"
	formatcache "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/format_cache"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/mappers"
)

type Core struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	formatCache *formatcache.FormatCache
	downloader  *downloader.Downloader

	// parameters
	ytDlpName      string
	ytDlpPath      string
	downloadsDir   string // Directory where downloaded files are saved
	serviceOptions *dto.Options
}

func NewCore(
	logger *slog.Logger,
	ytDlpPath string,
	downloadsDir string,
	serviceOptions *dto.Options,
) *Core {
	formatCache := formatcache.NewFormatCache(path.Join(downloadsDir, consts.YtDlpFormatCacheDir))

	return &Core{
		logger:  logger,
		mappers: mappers.NewMappers(),

		formatCache: formatCache,
		downloader:  downloader.NewDownloader(logger, formatCache, ytDlpPath, downloadsDir, serviceOptions),

		// parameters
		ytDlpName:      filepath.Base(ytDlpPath),
		ytDlpPath:      ytDlpPath,
		downloadsDir:   downloadsDir,
		serviceOptions: serviceOptions,
	}
}
