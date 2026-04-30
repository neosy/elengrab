package ffmpegsrv

import (
	"context"

	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// GetVideoAudioInfoFromFile extracts video and audio information from a media file
// using ffmpeg. Returns VideoInfo and AudioInfo, or nil if unavailable.
func (srv *FFmpegService) GetVideoAudioInfoFromFile(
	ctx context.Context,
	filePath string,
	srcMediaInfo *dservices.MediaInfo,
) (*dtypes.VideoInfo, *dtypes.AudioInfo, error) {
	return srv.core.GetVideoAudioInfo(ctx, filePath, srcMediaInfo)
}
