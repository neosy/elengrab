package thumbnail

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

// GetInfoByThumbID retrieves thumbnail information by its ID.
func (t *Thumbnail) GetByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error) {
	thumbnail, err := t.store.GetByThumbID(ctx, thumbID)
	if err != nil {
		return nil, err
	}

	if thumbnail.ImageRaw == nil {
		raw, err := t.storage.Get(thumbnail.StorageKey, thumbnail.Variant, thumbnail.Format.String())
		if err != nil {
			t.logger.Warn(
				"Failed to get thumbnail from storage",
				"error", err,
			)
			return nil, errorx.Errorf("failed to get thumbnail from storage: %w", err, exceptionx.NOT_FOUND)
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
