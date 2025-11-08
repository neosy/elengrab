package ytdlpsrv

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/neosy/elengrab/internal/services/ytdlp/mappers"
)

const (
	ytDlpName = "yt-dlp"
)

type YtDlpService struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	cmdPath string
	// Directory where downloaded files are saved
	downloadsDir string
}

func NewYtDlpService(logger *slog.Logger, binDir string, downloadsDir string) (*YtDlpService, error) {
	cmdPath, err := resolveCmdPath(ytDlpName, binDir)
	if err != nil {
		return nil, err
	}

	if downloadsDir == "" {
		return nil, errors.New("download directory is not set")
	}

	if strings.HasSuffix(downloadsDir, "/") || strings.HasSuffix(downloadsDir, "\\") {
		return nil, fmt.Errorf("downloads directory must not end with a slash or backslash: %s", downloadsDir)
	}

	if err := checkDir(downloadsDir); err != nil {
		return nil, err
	}

	return &YtDlpService{
		logger:  logger,
		mappers: mappers.NewMappers(),

		cmdPath:      cmdPath,
		downloadsDir: downloadsDir,
	}, nil
}
