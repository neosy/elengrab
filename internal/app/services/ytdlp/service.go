package ytdlpsrv

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/utils"
	"github.com/neosy/elengrab/pkg/nfile"
)

// Options holds configuration options for YtDlpService
type Options struct {
	// Number of fragments of a dash/hlsnative video that should be downloaded concurrently (default is 1)
	ConcurrentFragments uint8
}

// setDefaults sets default values for Options fields if they are not set
// or if force is true
func (o *Options) setDefaults(force bool) {
	if o.ConcurrentFragments == 0 || force {
		o.ConcurrentFragments = concurrentFragmentsDefault
	}
}

type YtDlpService struct {
	logger *slog.Logger

	// options
	options Options

	// internal
	core *core.Core
}

func NewYtDlpService(logger *slog.Logger, binDir string, downloadsDir string, options *Options) (*YtDlpService, error) {
	cmdPath, err := utils.ResolveCmdPath(ytDlpName, binDir)
	if err != nil {
		return nil, err
	}

	if err := utils.CheckCmd(ffmpegName); err != nil {
		return nil, err
	}

	if downloadsDir == "" {
		return nil, errors.New("download directory is not set")
	}

	if strings.HasSuffix(downloadsDir, "/") || strings.HasSuffix(downloadsDir, "\\") {
		return nil, fmt.Errorf("downloads directory must not end with a slash or backslash: %s", downloadsDir)
	}

	if err := nfile.CheckDir(downloadsDir); err != nil {
		return nil, err
	}

	var opts Options
	if options != nil {
		opts = *options
	}
	opts.setDefaults(false)

	return &YtDlpService{
		logger: logger,

		// options
		options: opts,

		// internal
		core: core.NewCore(logger, cmdPath, downloadsDir),
	}, nil
}
