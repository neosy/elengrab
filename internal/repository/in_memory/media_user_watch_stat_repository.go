package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type mediaUserWatchStatKey struct {
	downloadID uuid.UUID
	userID     uuid.UUID
}

type MediaUserWatchStatRepository struct {
	memsimple.Repository[ddownload.MediaUserWatchStat]

	cache memsimple.Cache[mediaUserWatchStatKey, ddownload.MediaUserWatchStat]
}

// newMediaUserWatchStatRepository returns a new object for the repository
func newMediaUserWatchStatRepository(ttl time.Duration) *MediaUserWatchStatRepository {
	r := &MediaUserWatchStatRepository{
		cache: memsimple.NewCacheWithDeaultCopier[mediaUserWatchStatKey, ddownload.MediaUserWatchStat, *ddownload.MediaUserWatchStat](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *MediaUserWatchStatRepository) Name() string {
	return "media_watch_stat"
}

func (r *MediaUserWatchStatRepository) Save(ctx context.Context, stat *ddownload.MediaUserWatchStat) error {
	if stat == nil {
		return nil
	}

	key := mediaUserWatchStatKey{
		downloadID: stat.DownloadID,
		userID:     stat.UserID,
	}

	save := func() error {
		r.cache.Save(
			key,
			stat,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *MediaUserWatchStatRepository) SaveNegative(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) error {
	if downloadID == uuid.Nil {
		return nil
	}

	key := mediaUserWatchStatKey{
		downloadID: downloadID,
		userID:     userID,
	}

	save := func() error {
		r.cache.Save(
			key,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

// Delete removes a mediaDownload from the in-memory repository using its ID.
func (r *MediaUserWatchStatRepository) Delete(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) error {
	key := mediaUserWatchStatKey{
		downloadID: downloadID,
		userID:     userID,
	}

	delete := func() error {
		if downloadID != uuid.Nil {
			r.cache.Delete(key)
		}
		return nil
	}
	return r.Repository.Delete(ctx, delete)
}

// FindByDownloadID retrieves a mediaDownload by its downloadID, userID from the repository.
func (r *MediaUserWatchStatRepository) Find(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
) (*ddownload.MediaUserWatchStat, memsimple.CacheStatus, error) {
	key := mediaUserWatchStatKey{
		downloadID: downloadID,
		userID:     userID,
	}

	find := func() (*ddownload.MediaUserWatchStat, memsimple.CacheStatus, error) {
		stat, status := r.cache.FindWithStatus(key)
		return stat, status, nil
	}

	return r.Repository.FindWithStatus(ctx, find)
}

// Exists checks if a mediaDownload exists in the repository by its downloadID, userID.
func (r *MediaUserWatchStatRepository) Exists(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
) (bool, memsimple.CacheStatus, error) {
	key := mediaUserWatchStatKey{
		downloadID: downloadID,
		userID:     userID,
	}

	exists := func() (bool, memsimple.CacheStatus, error) {
		exists, status := r.cache.ExistsWithStatus(key)
		return exists, status, nil
	}

	return r.Repository.ExistsWithStatus(ctx, exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *MediaUserWatchStatRepository) CleanExpired(ctx context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cache.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(ctx, clean)
}
