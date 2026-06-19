package repository

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (r *ThumbnailRepository) Insert(ctx context.Context, thumbnail *dmedia.Thumbnail) error {
	if thumbnail == nil {
		r.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	err := r.repo.Insert(ctx, thumbnail)
	if err != nil {
		r.logger.Warn(
			"Failed to insert record into repository",
			"thumbnailID", thumbnail.ThumbID,
			"storageKey", thumbnail.StorageKey,
			"sourceType", thumbnail.SourceType,
			"sourceURL", uptr.Deref(thumbnail.SourceURL),
			"variant", thumbnail.Variant,
			"format", thumbnail.Format.String(),
			"error", err,
		)
		return err
	}

	thumbnail, _ = r.repo.FindByThumbID(ctx, thumbnail.ThumbID)
	if thumbnail != nil {
		r.cacheRepo.Save(ctx, thumbnail)
	}

	return nil
}
