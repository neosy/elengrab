package ffmpegsrv

import (
	"context"
)

func (srv *FFmpegService) GetFFmpegVersion(ctx context.Context) (string, error) {
	return srv.core.GetFFmpegVersion(ctx)
}
