package sourceindex

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaSourceIndex) Save(ctx context.Context, index *ddownload.MediaSourceIndex) error {
	if index == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	if err := index.Validate(); err != nil {
		return err
	}

	err := uc.indexRepo().Save(ctx, index)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return errorx.Errorf("failed to insert record: %w", err, exceptionx.ERROR)
	}

	return nil

}

func (uc *MediaSourceIndex) Insert(ctx context.Context, index *ddownload.MediaSourceIndex) error {
	return uc.Save(ctx, index)
}

func (uc *MediaSourceIndex) Update(ctx context.Context, index *ddownload.MediaSourceIndex) error {
	return uc.Save(ctx, index)
}
