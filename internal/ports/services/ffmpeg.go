package pservices

import (
	"context"

	ffmpegsrv "github.com/neosy/elengrab/internal/app/services/ffmpeg"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type FFMpeg interface {
	ExtractBestFrame(
		ctx context.Context,
		filePath string,
		opts ...ffmpegsrv.FrameOption,
	) (*dtypes.ImageData, error)

	// GetVideoAudioInfoFromFile extracts video and audio information from a media file
	// using ffmpeg. Returns VideoInfo and AudioInfo, or nil if unavailable.
	GetVideoAudioInfoFromFile(
		ctx context.Context,
		filePath string,
		srcMediaInfo *dservices.MediaInfo,
	) (*dtypes.VideoInfo, *dtypes.AudioInfo, error)
}
