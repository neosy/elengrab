package linkclick

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *LinkClick) Create(ctx context.Context, linkClick *dlink.LinkClick) (uuid.UUID, error) {
	if linkClick == nil {
		uc.logger.Warn("Nil pointer in function")
		return uuid.Nil, apperrors.ErrFuncParamNullPointer
	}

	if linkClick.LinkClickID == uuid.Nil {
		linkClick.LinkClickID = uuid.New()
	}

	err := uc.linkClickRepo().Insert(ctx, linkClick)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return uuid.Nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return linkClick.LinkClickID, nil
}
