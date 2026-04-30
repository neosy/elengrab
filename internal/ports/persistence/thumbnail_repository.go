package persistence

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type ThumbnailRepository interface {
	Transactional
	Insert(ctx context.Context, thumbnail *dmedia.Thumbnail) error
	Delete(ctx context.Context, thumbID uuid.UUID) error
	FindByThumbID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error)
	ExistsByThumbID(ctx context.Context, thumbID uuid.UUID) (bool, error)
	FindByMediaIDBest(ctx context.Context, mediaID uuid.UUID) (*dmedia.Thumbnail, error)
	FindByVersion(
		ctx context.Context,
		mediaID uuid.UUID,
		variant dtypes.ThumbnailVariant,
		version uint8,
	) (*dmedia.Thumbnail, error)
	GetByMediaID(ctx context.Context, mediaID uuid.UUID) ([]*dmedia.Thumbnail, error)
}
