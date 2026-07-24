package watchevent

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaWatchEvent) Create(ctx context.Context, event *ddownload.MediaWatchEvent) error {
	if event == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	if event.EventID == uuid.Nil {
		event.EventID = uuid.New()
	}

	if err := event.Validate(); err != nil {
		return err
	}

	err := uc.eventRep.Insert(ctx, event)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return errorx.Errorf("failed to insert record: %w", err, exceptionx.ERROR)
	}

	return nil
}
