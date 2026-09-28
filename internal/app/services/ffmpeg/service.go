package ffmpegsrv

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/services/ffmpeg/internal/consts"
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
	cmdFFmpegPath, err := utils.ResolveCmdPath(consts.FFmpegName, binDir)
	if err != nil {
		return nil, err
	}

	cmdFFprobePath, err := utils.ResolveCmdPath(consts.FFprobeName, binDir)
	if err != nil {
		return nil, err
	}

	err = utils.CheckFFmpeg(consts.FFmpegName)
	if err != nil {
		return nil, err
	} else {
		logger.Info("FFmpeg executable found in PATH", "executable", consts.FFmpegName)
	}

	err = utils.CheckFFprobe(consts.FFprobeName)
	if err != nil {
		return nil, err
	} else {
		logger.Info("FFprobe executable found in PATH", "executable", consts.FFprobeName)
	}

	return &FFmpegService{
		logger: logger,
		core:   core.NewFFmpegCore(logger, cmdFFmpegPath, cmdFFprobePath),
	}, nil
}
