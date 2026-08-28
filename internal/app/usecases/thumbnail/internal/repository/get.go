package repository

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	ierrors "github.com/neosy/elengrab/internal/errors"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (r *ThumbnailRepository) FindByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error) {
	if thumbID == uuid.Nil {
		return nil, nil
	}

	thumbnail, cacheStatus, _ := r.cacheRepo.FindByThumbnailID(ctx, thumbID)
	if thumbnail != nil {
		return thumbnail, nil
	}
	if cacheStatus == memsimple.CacheStatusNegativeHit {
		return nil, nil
	}

	thumbnail, err := r.repo().FindByThumbID(ctx, thumbID)
	if err != nil {
		r.logger.Warn("Failed get thumbnail", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if thumbnail != nil {
		r.cacheRepo.Save(ctx, thumbnail)
	} else {
		r.cacheRepo.SaveNegative(ctx, thumbID)
	}

	return thumbnail, nil
}

func (r *ThumbnailRepository) GetByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error) {
	thumbnail, err := r.FindByThumbID(ctx, thumbID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if thumbnail == nil {
		r.logger.Warn("Thumbnail not found", "thumbID", thumbID)
		return nil, ierrors.ErrThumbnailNotFound
	}

	return thumbnail, nil
}
