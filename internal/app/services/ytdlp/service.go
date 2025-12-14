package ytdlpsrv

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	iutils "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/utils"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/ytdlp"
)

const (
	ytDlpName = "yt-dlp"
)

type Options struct {
	// Number of fragments of a dash/hlsnative video that should be downloaded concurrently (default is 1)
	ConcurrentFragments uint8
}

type YtDlpService struct {
	logger *slog.Logger

	// options
	options Options

	// internal
	ytdlp *ytdlp.YTDlp
}

func NewYtDlpService(logger *slog.Logger, binDir string, downloadsDir string, options *Options) (*YtDlpService, error) {
	cmdPath, err := iutils.ResolveCmdPath(ytDlpName, binDir)
	if err != nil {
		return nil, err
	}

	if downloadsDir == "" {
		return nil, errors.New("download directory is not set")
	}

	if strings.HasSuffix(downloadsDir, "/") || strings.HasSuffix(downloadsDir, "\\") {
		return nil, fmt.Errorf("downloads directory must not end with a slash or backslash: %s", downloadsDir)
	}

	if err := iutils.CheckDir(downloadsDir); err != nil {
		return nil, err
	}

	var opts Options
	if options != nil {
		opts = *options
	}

	return &YtDlpService{
		logger: logger,

		// options
		options: opts,

		// internal
		ytdlp: ytdlp.NewYTDlp(logger, cmdPath, downloadsDir),
	}, nil
}
