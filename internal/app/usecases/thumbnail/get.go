package thumbnail

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

// GetInfoByThumbID retrieves thumbnail information by its ID.
func (t *Thumbnail) GetByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error) {
	thumbnail, err := t.store.GetByThumbID(ctx, thumbID)
	if err != nil {
		return nil, err
	}

	getThumbnailFile := func() ([]byte, error) {
		fileID := dmedia.MakeThumbnailStorageKey(thumbID, thumbnail.Variant.String(), thumbnail.Format.String())
		raw, cacheStatus, err := t.thumbnailFileCache.FindByFileID(fileID)
		if err != nil {
			return nil, err
		}
		if cacheStatus == memsimple.CacheStatusHit {
			return raw, nil
		}
		if cacheStatus == memsimple.CacheStatusNegativeHit {
			t.logger.Warn(
				"Thumbnail file not found in cache (negative cache)",
				"thumbnailID", thumbID,
			)
			return nil, errorx.Errorf("thumbnail file not found (negative cache)", exceptionx.NOT_FOUND)
		}

		raw, err = t.storage.Get(thumbnail.StorageKey, thumbnail.Variant, thumbnail.Format.String())
		if err != nil {
			t.logger.Warn(
				"Failed to get thumbnail from storage",
				"thumbnailID", thumbID,
				"error", err,
			)
			t.thumbnailFileCache.SaveNegative(fileID)
			return nil, errorx.Errorf("failed to get thumbnail from storage: %w", err, exceptionx.NOT_FOUND)
		}

		if len(raw) == 0 {
			t.thumbnailFileCache.SaveNegative(fileID)
			return nil, errorx.Errorf("thumbnail file not found", exceptionx.NOT_FOUND)
		}

		t.thumbnailFileCache.Save(fileID, raw)

		return raw, nil
	}

	if thumbnail.ImageRaw == nil {
		raw, err := getThumbnailFile()
		if err != nil {
			return nil, err
		}
		thumbnail.ImageRaw = raw
	}

	return thumbnail, nil
}

// GetInfoByThumbID retrieves thumbnail information by its ID without fetching the raw image data.
func (t *Thumbnail) GetInfoByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error) {
	thumbnail, err := t.store.GetByThumbID(ctx, thumbID)
	if err != nil {
		return nil, err
	}

	return thumbnail, nil
}

func (t *Thumbnail) FindInfoByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error) {
	thumbnail, err := t.store.FindByThumbID(ctx, thumbID)
	if err != nil {
		return nil, err
	}

	return thumbnail, nil
}
