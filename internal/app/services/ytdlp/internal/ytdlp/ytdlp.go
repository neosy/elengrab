package ytdlp

import (
	"log/slog"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/ytdlp/mappers"
)

type YTDlp struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// parameters
	ytDlpName    string
	ytDlpPath    string
	downloadsDir string // Directory where downloaded files are saved
}

func NewYTDlp(logger *slog.Logger, ytDlpPath string, downloadsDir string) *YTDlp {
	return &YTDlp{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// parameters
		ytDlpName:    filepath.Base(ytDlpPath),
		ytDlpPath:    ytDlpPath,
		downloadsDir: downloadsDir,
	}
}
