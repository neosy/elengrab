package pservices

import (
	"context"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type FFMpeg interface {
	GetFFmpegVersion(ctx context.Context) (string, error)

	ExtractBestFrame(
		ctx context.Context,
		filePath string,
		opts ...FrameOption,
	) (*dtypes.ImageData, error)

	// ExtractMediaInfoWithFFprobe extracts video and audio information from a media file
	// using ffmpeg. Returns VideoInfo and AudioInfo, or nil if unavailable.
	ExtractMediaInfoWithFFprobe(
		ctx context.Context,
		filePath string,
		srcMediaInfo *MediaInfo,
	) (*MediaInfo, error)

	// ExtractDurationMs extracts media duration from the given file using ffprobe.
	//
	// It returns duration in milliseconds. If ffprobe fails or duration cannot be parsed,
	// an error is returned.
	ExtractDurationMs(ctx context.Context, filePath string) (time.Duration, error)
}
