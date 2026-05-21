package ffmpegsrv

import (
	"context"

	dservices "github.com/neosy/elengrab/internal/domain/services"
)

// ExtractVideoAudioInfoFromFile extracts video and audio information from a media file
// using ffmpeg. Returns VideoInfo and AudioInfo, or nil if unavailable.
func (srv *FFmpegService) ExtractVideoAudioInfoFromFile(
	ctx context.Context,
	filePath string,
	srcMediaInfo *dservices.MediaInfo,
) (*dservices.MediaInfo, error) {
	return srv.core.ExtractVideoAudioInfoWithFFprobe(ctx, filePath, srcMediaInfo)
}
