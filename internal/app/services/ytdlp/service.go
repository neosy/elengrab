package ytdlpsrv

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/mappers"
)

const (
	ytDlpName                  = "yt-dlp"
	concurrentFragmentsDefault = 5
)

type Options struct {
	// Number of fragments of a dash/hlsnative video that should be downloaded concurrently (default is 1)
	ConcurrentFragments uint8
}

type YtDlpService struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	cmdPath string
	// Directory where downloaded files are saved
	downloadsDir string

	options Options
}

func NewYtDlpService(logger *slog.Logger, binDir string, downloadsDir string, options *Options) (*YtDlpService, error) {
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

	var opts Options
	if options != nil {
		opts = *options
	}

	if opts.ConcurrentFragments == 0 || opts.ConcurrentFragments > 20 {
		opts.ConcurrentFragments = concurrentFragmentsDefault
	}

	return &YtDlpService{
		logger:  logger,
		mappers: mappers.NewMappers(),

		cmdPath:      cmdPath,
		downloadsDir: downloadsDir,

		options: opts,
	}, nil
}
