package downloader

import (
	"log/slog"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	formatcache "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/format_cache"
)

type Downloader struct {
	logger      *slog.Logger
	formatCache *formatcache.FormatCache

	// parameters
	ytDlpName      string
	ytDlpPath      string
	downloadsDir   string // Directory where downloaded files are saved
	serviceOptions *dto.Options
}

// NewDownloader creates a new Downloader instance.
func NewDownloader(
	logger *slog.Logger,
	formatCache *formatcache.FormatCache,
	ytDlpPath string,
	downloadsDir string,
	serviceOptions *dto.Options,
) *Downloader {
	return &Downloader{
		logger:      logger,
		formatCache: formatCache,

		// parameters
		ytDlpName:      filepath.Base(ytDlpPath),
		ytDlpPath:      ytDlpPath,
		downloadsDir:   downloadsDir,
		serviceOptions: serviceOptions,
	}
}
