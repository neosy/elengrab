package ffmpegsrv

import (
	"context"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

// ExtractMediaInfoWithFFprobe extracts video and audio information from a media file
// using ffmpeg. Returns VideoInfo and AudioInfo, or nil if unavailable.
func (srv *FFmpegService) ExtractMediaInfoWithFFprobe(
	ctx context.Context,
	filePath string,
	baseMediaInfo *pservices.MediaInfo,
) (*pservices.MediaInfo, error) {
	return srv.core.ExtractMediaInfoWithFFprobe(ctx, filePath, baseMediaInfo)
}
