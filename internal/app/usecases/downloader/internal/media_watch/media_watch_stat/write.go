package watchstat

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaWatchStat) Write(ctx context.Context, stat *ddownload.MediaWatchStat) error {
	if stat == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	if err := stat.Validate(); err != nil {
		return err
	}

	err := uc.statRep.Write(ctx, stat)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return errorx.Errorf("failed to insert record: %w", err, exceptionx.ERROR)
	}

	uc.statCacheRep.Save(ctx, stat)

	return nil
}
