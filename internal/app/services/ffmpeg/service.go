package ffmpegsrv

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/services/ffmpeg/internal/core"
	"github.com/neosy/elengrab/internal/app/services/ffmpeg/internal/utils"
)

// FFmpegService represents a service for interacting with ffmpeg.
type FFmpegService struct {
	logger *slog.Logger

	// internal
	core *core.FFmpegCore
}

// NewFFmpegService creates a new instance of FfmpegService.
func NewFFmpegService(
	logger *slog.Logger,
	binDir string,
) (*FFmpegService, error) {
	cmdPath, err := utils.ResolveCmdPath(ffmpegName, binDir)
	if err != nil {
		return nil, err
	}

	err = utils.CheckFFmpeg(ffmpegName)
	if err != nil {
		return nil, err
	} else {
		logger.Debug("FFmpeg executable found in PATH", "executable", ffmpegName)
	}

	return &FFmpegService{
		logger: logger,
		core:   core.NewFFmpegCore(logger, cmdPath),
	}, nil
}
