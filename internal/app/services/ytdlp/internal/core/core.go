package core

import (
	"log/slog"
	"path"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/downloader"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/mappers"
)

type Core struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	formatCache *idto.FormatCache
	downloader  *downloader.Downloader

	// parameters
	ytDlpName    string
	ytDlpPath    string
	downloadsDir string // Directory where downloaded files are saved
}

func NewCore(logger *slog.Logger, ytDlpPath string, downloadsDir string) *Core {

	formatCache := idto.NewFormatCache(path.Join(downloadsDir, ytDlpFormatCacheDir))

	return &Core{
		logger:  logger,
		mappers: mappers.NewMappers(),

		formatCache: formatCache,
		downloader:  downloader.NewDownloader(logger, formatCache, ytDlpPath, downloadsDir),

		// parameters
		ytDlpName:    filepath.Base(ytDlpPath),
		ytDlpPath:    ytDlpPath,
		downloadsDir: downloadsDir,
	}
}
