package ffmpegsrv

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// ExtractBestFrame extracts the best frame from a media file using ffmpeg. Returns the frame as a byte slice.
func (srv *FFmpegService) ExtractBestFrame(
	ctx context.Context,
	filePath string,
	opts ...FrameOption,
) (*dtypes.ImageData, error) {
	options := FrameOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.Strategy == nil {
		options.Strategy = FrameStrategyThumbnail{}
	}

	if options.Format == nil {
		options.Format = FrameFormatJPEG{}
	}

	args := options.Strategy.Args()
	args = append(args, options.Format.Args()...)

	imgData, err := srv.core.GetFrame(ctx, filePath, args...)
	if err != nil {
		return nil, err
	}

	imgData.Format = options.Format.Format()

	return imgData, nil
}
