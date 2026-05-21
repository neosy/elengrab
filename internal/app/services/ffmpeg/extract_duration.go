package ffmpegsrv

import (
	"context"
	"time"
)

// ExtractDurationMs extracts media duration from the given file using ffprobe.
//
// It returns duration in milliseconds. If ffprobe fails or duration cannot be parsed,
// an error is returned.
func (srv *FFmpegService) ExtractDurationMs(ctx context.Context, filePath string) (time.Duration, error) {
	return srv.core.ExtractDurationMsWithFFprobe(ctx, filePath)
}
