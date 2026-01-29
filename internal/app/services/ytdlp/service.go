package ytdlpsrv

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core"
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
	cmdPath, err := resolveCmdPath(ytDlpName, binDir)
	if err != nil {
		return nil, err
	}

	if err := checkFFmpeg(); err != nil {
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

func resolveCmdPath(cmdName, binDir string) (string, error) {
	// On Windows, add .exe suffix if missing
	if runtime.GOOS == "windows" && !strings.HasSuffix(cmdName, ".exe") {
		cmdName += ".exe"
	}

	// try PATH
	if path, err := lookupExecutable(cmdName); err == nil {
		return path, nil
	}

	// try config dir
	cmdPath := filepath.Join(binDir, cmdName)
	if fi, err := os.Stat(cmdPath); err == nil && !fi.IsDir() {
		return cmdPath, nil
	}

	// 3) try executable directory (same folder as your service binary)
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		cmdPath := filepath.Join(exeDir, cmdName)
		if fi, err := os.Stat(cmdPath); err == nil && !fi.IsDir() {
			return cmdPath, nil
		}
	}

	return "", fmt.Errorf(
		"%q executable not found. Tried PATH lookup and config directory %q (full path: %q)",
		cmdName,
		binDir,
		cmdPath,
	)
}

func lookupExecutable(cmdName string) (string, error) {
	if runtime.GOOS == "windows" && !strings.HasSuffix(cmdName, ".exe") {
		cmdName += ".exe"
	}
	return exec.LookPath(cmdName)
}

func checkFFmpeg() error {
	var cmdName = "ffmpeg"

	if runtime.GOOS == "windows" && !strings.HasSuffix(cmdName, ".exe") {
		cmdName += ".exe"
	}

	// Ensure ffmpeg is available in PATH
	cmd := exec.Command(cmdName, "-version")
	if err := cmd.Run(); err == nil {
		return nil
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		cmdPath := filepath.Join(exeDir, cmdName)

		cmd = exec.Command(cmdPath, "-version")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return errors.New("ffmpeg not found in PATH. Please install ffmpeg and add it to PATH")
}
