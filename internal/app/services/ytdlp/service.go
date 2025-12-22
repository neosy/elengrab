package ytdlpsrv

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ytdlp "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/yt_dlp"
	"github.com/neosy/elengrab/pkg/nfile"
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

	if err := nfile.CheckDir(downloadsDir); err != nil {
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

func resolveCmdPath(cmdName, binDir string) (string, error) {
	if path, err := exec.LookPath(cmdName); err == nil {
		return path, nil
	}

	cmdPath := filepath.Join(binDir, cmdName)
	if fi, err := os.Stat(cmdPath); err == nil && !fi.IsDir() {
		return cmdPath, nil
	}

	return "", fmt.Errorf("%s not found: tried config path %q and PATH lookup", cmdName, cmdPath)
}
