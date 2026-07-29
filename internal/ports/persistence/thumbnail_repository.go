package persistence

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
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

type ThumbnailCacheRepository interface {
	memory.CacheRepository

	Save(ctx context.Context, thumbnail *dmedia.Thumbnail) error
	SaveNegative(ctx context.Context, thumbID uuid.UUID) error

	FindByThumbnailID(ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, memsimple.CacheStatus, error)
	ExistsByThumbnailID(ctx context.Context, thumbID uuid.UUID) (bool, error)

	CleanExpired(context.Context) error
}

type ThumbnailFileCacheRepository interface {
	memory.CacheRepository

	Save(thumbID uuid.UUID, file []byte) error
	SaveNegative(fileID uuid.UUID) error
	Delete(fileID uuid.UUID) error

	FindByFileID(fileID uuid.UUID) ([]byte, memsimple.CacheStatus, error)
	ExistsByFileID(fileID uuid.UUID) (bool, error)

	CleanExpired(context.Context) error
}
