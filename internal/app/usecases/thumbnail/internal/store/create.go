package store

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (c *ThumbnailStore) Insert(ctx context.Context, thumbnail *dmedia.Thumbnail) error {
	if thumbnail == nil {
		c.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	err := c.thumbnailRep.Insert(ctx, thumbnail)
	if err != nil {
		c.logger.Warn(
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

	return nil
}
