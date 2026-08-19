package uwatchposition

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaUserWatchPosition) Write(ctx context.Context, position *ddownload.MediaUserWatchPosition) error {
	if position == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	if err := position.Validate(); err != nil {
		uc.logger.Debug(
			"Media watch position validation failed",
			"error", err,
		)
		return err
	}

	err := uc.positionRep.Write(ctx, position)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return errorx.Errorf("failed to insert record: %w", err, exceptionx.ERROR)
	}

	uc.positionCacheRep.Save(ctx, position)

	return nil
}
