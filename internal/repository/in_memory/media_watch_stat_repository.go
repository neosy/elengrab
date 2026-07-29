package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type MediaWatchStatRepository struct {
	memsimple.Repository[ddownload.MediaWatchStat]

	cacheByDownloadID memsimple.Cache[uuid.UUID, ddownload.MediaWatchStat]
}

// newMediaWatchStatRepository returns a new object for the repository
func newMediaWatchStatRepository(ttl time.Duration) *MediaWatchStatRepository {
	r := &MediaWatchStatRepository{
		cacheByDownloadID: memsimple.NewCacheWithDeaultCopier[uuid.UUID, ddownload.MediaWatchStat, *ddownload.MediaWatchStat](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *MediaWatchStatRepository) Name() string {
	return "media_watch_stat"
}

func (r *MediaWatchStatRepository) Save(stat *ddownload.MediaWatchStat) error {
	if stat == nil {
		return nil
	}

	save := func() error {
		r.cacheByDownloadID.Save(
			stat.DownloadID,
			stat,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

func (r *MediaWatchStatRepository) SaveNegative(downloadID uuid.UUID) error {
	if downloadID == uuid.Nil {
		return nil
	}

	save := func() error {
		r.cacheByDownloadID.Save(
			downloadID,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

// Delete removes a mediaDownload from the in-memory repository using its ID.
func (r *MediaWatchStatRepository) Delete(downloadID uuid.UUID) error {
	delete := func() error {
		if downloadID != uuid.Nil {
			r.cacheByDownloadID.Delete(downloadID)
		}
		return nil
	}
	return r.Repository.Delete(delete)
}

// FindByDownloadID retrieves a mediaDownload by its downloadID from the repository.
func (r *MediaWatchStatRepository) Find(downloadID uuid.UUID) (*ddownload.MediaWatchStat, memsimple.CacheStatus, error) {
	find := func() (*ddownload.MediaWatchStat, memsimple.CacheStatus, error) {
		stat, status := r.cacheByDownloadID.FindWithStatus(downloadID)
		return stat, status, nil
	}

	return r.Repository.FindWithStatus(find)
}

// Checks if a mediaDownload exists by its downloadID.
func (r *MediaWatchStatRepository) Exists(downloadID uuid.UUID) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByDownloadID.Exists(downloadID), nil
	}
	return r.Repository.Exists(exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *MediaWatchStatRepository) CleanExpired(context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByDownloadID.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(clean)
}
