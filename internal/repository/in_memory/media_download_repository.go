package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type MediaDownloadRepository struct {
	memsimple.Repository[ddownload.MediaDownload]

	cacheByDownloadID memsimple.Cache[uuid.UUID, ddownload.MediaDownload]
}

// newMediaDownloadRepository returns a new object for the repository
func newMediaDownloadRepository(ttl time.Duration) *MediaDownloadRepository {
	r := &MediaDownloadRepository{
		cacheByDownloadID: memsimple.NewCacheWithDeaultCopier[uuid.UUID, ddownload.MediaDownload, *ddownload.MediaDownload](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *MediaDownloadRepository) Name() string {
	return "media_download"
}

func (r *MediaDownloadRepository) Save(ctx context.Context, media *ddownload.MediaDownload) error {
	if media == nil {
		return nil
	}

	save := func() error {
		r.cacheByDownloadID.Save(
			media.DownloadID,
			media,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *MediaDownloadRepository) SaveNegative(ctx context.Context, downloadID uuid.UUID) error {
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

	return r.Repository.Save(ctx, save)
}

// Delete removes a mediaDownload from the in-memory repository using its ID.
func (r *MediaDownloadRepository) Delete(ctx context.Context, downloadID uuid.UUID) error {
	delete := func() error {
		if downloadID != uuid.Nil {
			r.cacheByDownloadID.Delete(downloadID)
		}
		return nil
	}
	return r.Repository.Delete(ctx, delete)
}

// FindByDownloadID retrieves a mediaDownload by its downloadID from the repository.
func (r *MediaDownloadRepository) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaDownload, memsimple.CacheStatus, error) {
	find := func() (*ddownload.MediaDownload, memsimple.CacheStatus, error) {
		media, status := r.cacheByDownloadID.FindWithStatus(downloadID)
		return media, status, nil
	}

	return r.Repository.FindWithStatus(ctx, find)
}

// Checks if a mediaDownload exists by its downloadID.
func (r *MediaDownloadRepository) ExistsByFileID(ctx context.Context, downloadID uuid.UUID) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByDownloadID.Exists(downloadID), nil
	}
	return r.Repository.Exists(ctx, exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *MediaDownloadRepository) CleanExpired(ctx context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByDownloadID.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(ctx, clean)
}
