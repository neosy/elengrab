package migrations

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type fetchThumbnail func(ctx context.Context, path string) *dtypes.ImageData

var makeFetchThumbnail = func(
	logger *slog.Logger,
	retryCount int,
	retryDelay time.Duration,
	run func(ctx context.Context, path string) (*dtypes.ImageData, error),
) fetchThumbnail {
	return func(ctx context.Context, path string) *dtypes.ImageData {
		if retryCount < 2 {
			imageData, _ := run(ctx, path)
			return imageData
		}

		timer := time.NewTimer(0)
		timerStop := func() {
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
		}
		defer timerStop()

		timerStop()

		for i := range retryCount {
			logger.Debug("Start fetch thumbnail", "attemp", i+1)
			imageData, err := run(ctx, path)
			if err == nil && imageData != nil {
				return imageData
			}

			timer.Reset(retryDelay)

			select {
			case <-ctx.Done():
				return nil
			case <-timer.C:
			}
		}

		return nil
	}
}

func (m *migrations) createThumbnail(
	ctx context.Context,
	imageData *dtypes.ImageData,
	sourceType dtypes.ThumbnailSourceType,
) (uuid.UUID, error) {
	if imageData == nil {
		return uuid.Nil, fmt.Errorf("image data is nil")
	}

	var sourceURL *string
	if imageData.URL != "" {
		sourceURL = &imageData.URL
	}

	req := &dto.CreateThumbnailRequest{
		SourceType: sourceType,
		SourceURL:  sourceURL,
		ImageData:  imageData,
	}

	return m.usecases.thumbnail.Create(ctx, req)
}
