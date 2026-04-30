package store

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (c *ThumbnailStore) FindByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error) {
	thumbnail, err := c.thumbnailRep.FindByThumbID(ctx, thumbID)
	if err != nil {
		c.logger.Warn("Failed get thumbnail", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}
	return thumbnail, nil
}

func (c *ThumbnailStore) GetByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error) {
	thumbnail, err := c.FindByThumbID(ctx, thumbID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if thumbnail == nil {
		c.logger.Warn("Thumbnail not found", "thumbID", thumbID)
		return nil, ierrors.ErrThumbnailNotFound
	}

	return thumbnail, nil
}
