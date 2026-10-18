package uploadexecutor

import (
	"context"

	"github.com/neosy/elengrab/internal/app/services/ffmpeg/ffframe"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (uc *Executor) extractThumbnailFromFile(ctx context.Context, filePath string) *dtypes.ImageData {
	extractBestFrame := func(ctx context.Context, filePath string) (*dtypes.ImageData, error) {
		imgData, err := uc.ffmpeg.ExtractBestFrame(
			ctx,
			filePath,
			pservices.WithFrameStrategy(ffframe.FrameStrategyBalanced{}),
			pservices.WithFrameFormat(ffframe.FrameFormatWebP{}),
		)
		if err != nil {
			imgData, err = uc.ffmpeg.ExtractBestFrame(ctx, filePath)
		}

		if err != nil {
			uc.logger.Warn(
				"Failed to get thumbnail from file",
				"filePath", filePath,
				"error", err,
			)
			return nil, err
		}

		return imgData, nil
	}

	var err error
	imgData, err := extractBestFrame(ctx, filePath)
	if err != nil {
		return nil
	}

	return imgData
}
