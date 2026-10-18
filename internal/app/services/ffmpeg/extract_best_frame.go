package ffmpegsrv

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

// ExtractBestFrame extracts the best frame from a media file using ffmpeg. Returns the frame as a byte slice.
func (srv *FFmpegService) ExtractBestFrame(
	ctx context.Context,
	filePath string,
	opts ...pservices.FrameOption,
) (*dtypes.ImageData, error) {
	options := newFrameOptionsDefault()
	for _, opt := range opts {
		opt(&options)
	}

	args := options.Strategy.Args()
	args = append(args, options.Format.Args()...)

	imgData, err := srv.core.ExtractFrame(ctx, filePath, args...)
	if err != nil {
		return nil, err
	}

	imgData.Format = options.Format.Format()

	return imgData, nil
}
