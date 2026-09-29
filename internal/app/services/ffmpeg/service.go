package ffmpegsrv

import (
	"log/slog"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/services/ffmpeg/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ffmpeg/internal/core"
	"github.com/neosy/elengrab/internal/app/services/internal/utils"
)

// FFmpegService represents a service for interacting with ffmpeg.
type FFmpegService struct {
	logger *slog.Logger

	// internal
	core *core.FFmpegCore
}

// NewFFmpegService creates a new instance of FFmpegService.
func NewFFmpegService(
	logger *slog.Logger,
	binDir string,
) (*FFmpegService, error) {
	ffmpegPath, err := utils.ResolveCmdPath(consts.FFmpegName, binDir)
	if err != nil {
		return nil, err
	}

	ffprobePath, err := utils.ResolveCmdPath(consts.FFprobeName, binDir)
	if err != nil {
		return nil, err
	}

	err = checkFFmpeg(ffmpegPath)
	if err != nil {
		return nil, err
	}

	logger.Info("FFmpeg executable found",
		"name", filepath.Base(ffmpegPath),
		"dir", filepath.Dir(ffmpegPath),
	)

	err = checkFFprobe(ffprobePath)
	if err != nil {
		return nil, err
	}

	logger.Info("FFprobe executable found",
		"name", filepath.Base(ffprobePath),
		"dir", filepath.Dir(ffprobePath),
	)

	return &FFmpegService{
		logger: logger,
		core:   core.NewFFmpegCore(logger, ffmpegPath, ffprobePath),
	}, nil
}
