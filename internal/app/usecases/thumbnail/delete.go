package thumbnail

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (t *thumbnail) Delete(ctx context.Context, thumbID uuid.UUID) error {
	thumbnail, err := t.repo.GetByThumbID(ctx, thumbID)
	if err != nil {
		return err
	}

	delete := func(ctx context.Context) error {
		err := t.repo.Delete(ctx, thumbnail.ThumbID)
		if err != nil {
			return err
		}

		err = t.storage.Delete(thumbnail.StorageKey, thumbnail.Variant, thumbnail.Format.String())
		if err != nil {
			t.logger.Warn(
				"Failed to delete thumbnail from storage",
				"thumbnailID", thumbID,
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

	err = t.repo.TxIndependent(ctx, delete)
	if err != nil {
		return err
	}

	fileID := dmedia.MakeThumbnailStorageKey(thumbID, thumbnail.Variant.String(), thumbnail.Format.String())
	t.thumbnailFileCache.Delete(ctx, fileID)

	t.logger.Debug(
		"Thumbnail deleted",
		"thumbnailID", thumbID,
		"storageKey", thumbnail.StorageKey,
		"sourceType", thumbnail.SourceType,
		"sourceURL", uptr.Deref(thumbnail.SourceURL),
		"variant", thumbnail.Variant,
		"format", thumbnail.Format.String(),
	)

	return nil
}

func (t *thumbnail) DeleteBatch(ctx context.Context, thumbIDs []uuid.UUID) error {
	for _, thumbID := range thumbIDs {
		err := t.Delete(ctx, thumbID)
		if err != nil {
			return err
		}
	}

	return nil
}
