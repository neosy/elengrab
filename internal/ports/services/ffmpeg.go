package pservices

import (
	"context"
	"time"

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

	// ExtractVideoAudioInfoFromFile extracts video and audio information from a media file
	// using ffmpeg. Returns VideoInfo and AudioInfo, or nil if unavailable.
	ExtractVideoAudioInfoFromFile(
		ctx context.Context,
		filePath string,
		srcMediaInfo *dservices.MediaInfo,
	) (*dservices.MediaInfo, error)

	// ExtractDurationMs extracts media duration from the given file using ffprobe.
	//
	// It returns duration in milliseconds. If ffprobe fails or duration cannot be parsed,
	// an error is returned.
	ExtractDurationMs(ctx context.Context, filePath string) (time.Duration, error)
}
