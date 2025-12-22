package ytdlp

import (
	"log/slog"
	"path"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/yt_dlp/mappers"
)

type YTDlp struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	formatCache *formatCache

	// parameters
	ytDlpName    string
	ytDlpPath    string
	downloadsDir string // Directory where downloaded files are saved
}

func NewYTDlp(logger *slog.Logger, ytDlpPath string, downloadsDir string) *YTDlp {
	return &YTDlp{
		logger:  logger,
		mappers: mappers.NewMappers(),

		formatCache: NewFormatCache(path.Join(downloadsDir, ytDlpFormatCacheDir)),

		// parameters
		ytDlpName:    filepath.Base(ytDlpPath),
		ytDlpPath:    ytDlpPath,
		downloadsDir: downloadsDir,
	}
}
