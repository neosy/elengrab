package uploadexecutor

import (
	"context"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *Executor) CreateThumbnail(
	ctx context.Context,
	imageDate *dtypes.ImageData,
	sourceType dtypes.ThumbnailSourceType,
) (uuid.UUID, error) {
	req := &dto.CreateThumbnailRequest{
		SourceType: sourceType,
		ImageData:  imageDate,
	}

	id, err := uc.thumbnail.Create(ctx, req)
	if err != nil {
		return uuid.Nil(), err
	}

	return id, nil
}
