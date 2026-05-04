package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

// Defines the structure for the in-memory repository.
type ThumbnailRepository struct {
	// Embeds the base Repository
	memsimple.Repository[dmedia.Thumbnail]

	// Cache for storing Thumbnails by their thumbnail ID.
	cacheByThumbnailID memsimple.Cache[uuid.UUID, dmedia.Thumbnail]
}

// newThumbnailRepository returns a new object for the repository
func newThumbnailRepository(ttl time.Duration) *ThumbnailRepository {
	r := &ThumbnailRepository{
		cacheByThumbnailID: memsimple.NewCacheWithDeaultCopier[uuid.UUID, dmedia.Thumbnail, *dmedia.Thumbnail](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *ThumbnailRepository) Save(_ context.Context, thumbnail *dmedia.Thumbnail) error {
	if thumbnail == nil || thumbnail.ThumbID == uuid.Nil {
		return nil
	}

	save := func() error {
		r.cacheByThumbnailID.Save(
			thumbnail.ThumbID,
			thumbnail,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

func (r *ThumbnailRepository) SaveNegative(_ context.Context, thumbID uuid.UUID) error {
	if thumbID == uuid.Nil {
		return nil
	}

	save := func() error {
		r.cacheByThumbnailID.Save(
			thumbID,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

// Delete removes a thumbnail from the in-memory repository using its ID.
func (r *ThumbnailRepository) Delete(_ context.Context, thumbID uuid.UUID) error {
	delete := func() error {
		if thumbID != uuid.Nil {
			r.cacheByThumbnailID.Delete(thumbID)
		}
		return nil
	}
	return r.Repository.Delete(delete)
}

// FindByThumbID retrieves a thumbnail by its thumbnail ID from the repository.
func (r *ThumbnailRepository) FindByThumbnailID(
	_ context.Context,
	thumbID uuid.UUID,
) (*dmedia.Thumbnail, memsimple.CacheStatus, error) {
	find := func() (*dmedia.Thumbnail, memsimple.CacheStatus, error) {
		data, status := r.cacheByThumbnailID.FindWithStatus(thumbID)
		return data, status, nil
	}
	return r.Repository.FindWithStatus(find)
}

// Checks if a thumbnail exists by its thumbnail ID.
func (r *ThumbnailRepository) ExistsByThumbnailID(_ context.Context, thumbID uuid.UUID) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByThumbnailID.Exists(thumbID), nil
	}
	return r.Repository.Exists(exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *ThumbnailRepository) CleanExpired(_ context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByThumbnailID.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(clean)
}
