package downloader

import (
	"log/slog"
	"path/filepath"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/dto"
)

type Downloader struct {
	logger      *slog.Logger
	formatCache *idto.FormatCache

	// parameters
	ytDlpName    string
	ytDlpPath    string
	downloadsDir string // Directory where downloaded files are saved
}

// NewDownloader creates a new Downloader instance.
func NewDownloader(
	logger *slog.Logger,
	formatCache *idto.FormatCache,
	ytDlpPath string,
	downloadsDir string,
) *Downloader {
	return &Downloader{
		logger:      logger,
		formatCache: formatCache,

		// parameters
		ytDlpName:    filepath.Base(ytDlpPath),
		ytDlpPath:    ytDlpPath,
		downloadsDir: downloadsDir,
	}
}
