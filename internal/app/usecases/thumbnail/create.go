package thumbnail

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (t *thumbnail) Create(ctx context.Context, req *dto.CreateThumbnailRequest) (uuid.UUID, error) {
	thumbnail, err := t.mappers.MapCreateThumbnailRequestToThumbnailDomain(req)
	if err != nil {
		return uuid.Nil, err
	}

	err = thumbnail.InitStorageKey()
	if err != nil {
		t.logger.Warn(
			"Failed to initialize storage key",
			"error", err,
		)
		return uuid.Nil, err
	}

	create := func(ctx context.Context) error {
		var err error
		err = t.repo.Insert(ctx, thumbnail)
		if err != nil {
			return err
		}

		err = t.storage.Put(req.ImageData.Raw, thumbnail.StorageKey, thumbnail.Variant, req.ImageData.Format.String())
		if err != nil {
			t.logger.Warn(
				"Failed to put thumbnail to storage",
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

	err = t.repo.Tx(ctx, create)
	if err != nil {
		return uuid.Nil, err
	}

	t.logger.Debug(
		"Thumbnail created",
		"thumbnailID", thumbnail.ThumbID,
		"storageKey", thumbnail.StorageKey,
		"sourceType", thumbnail.SourceType,
		"sourceURL", uptr.Deref(thumbnail.SourceURL),
		"variant", thumbnail.Variant,
		"format", thumbnail.Format.String(),
	)

	return thumbnail.ThumbID, nil
}
