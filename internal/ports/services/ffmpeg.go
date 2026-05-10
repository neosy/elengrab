package pservices

import (
	"context"

	ffmpegsrv "github.com/neosy/elengrab/internal/app/services/ffmpeg"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type FFMpeg interface {
	ExtractBestFrame(
		ctx context.Context,
		filePath string,
		opts ...ffmpegsrv.FrameOption,
	) (*dtypes.ImageData, error)
}
