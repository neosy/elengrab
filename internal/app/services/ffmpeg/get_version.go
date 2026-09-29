package ffmpegsrv

import (
	"context"

	"github.com/neosy/elengrab/internal/app/services/ffmpeg/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ffmpeg/internal/core"
	"github.com/neosy/elengrab/internal/app/services/internal/utils"
)

func (srv *FFmpegService) GetFFmpegVersion(ctx context.Context) (string, error) {
	return srv.core.GetFFmpegVersion(ctx)
}

// GetFFmpegVersion retrieves the version of the FFmpeg executable.
func GetFFmpegVersion(ctx context.Context) (string, error) {
	cmdPath, err := utils.ResolveCmdPath(consts.FFmpegName, "")
	if err != nil {
		return "", err
	}

	return core.GetFFmpegVersion(ctx, cmdPath)
}
