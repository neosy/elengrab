package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type MediaUserWatchPositionRepository struct {
	memsimple.Repository[ddownload.MediaUserWatchPosition]

	cache memsimple.Cache[string, ddownload.MediaUserWatchPosition]
}

func buildKey(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) string {
	const sep = "_"

	if downloadID == uuid.Nil {
		return ""
	}

	if userID == uuid.Nil && (sessionID == nil || *sessionID == uuid.Nil) {
		return ""
	}

	key := downloadID.String() + sep + userID.String()

	if userID == uuid.Nil && sessionID != nil {
		key += sep + sessionID.String()
	}

	return key
}

// newMediaUserWatchPositionRepository returns a new object for the repository
func newMediaUserWatchPositionRepository(ttl time.Duration) *MediaUserWatchPositionRepository {
	r := &MediaUserWatchPositionRepository{
		cache: memsimple.NewCacheWithDeaultCopier[string, ddownload.MediaUserWatchPosition, *ddownload.MediaUserWatchPosition](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *MediaUserWatchPositionRepository) Save(position *ddownload.MediaUserWatchPosition) error {
	if position == nil {
		return nil
	}

	key := buildKey(position.DownloadID, position.UserID, position.SessionID)
	if key == "" {
		return nil
	}

	save := func() error {
		r.cache.Save(
			key,
			position,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

func (r *MediaUserWatchPositionRepository) SaveNegative(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) error {
	if downloadID == uuid.Nil {
		return nil
	}

	key := buildKey(downloadID, userID, sessionID)
	if key == "" {
		return nil
	}

	save := func() error {
		r.cache.Save(
			key,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

// Delete removes a mediaDownload from the in-memory repository using its ID.
func (r *MediaUserWatchPositionRepository) Delete(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) error {
	fn := func() error {
		if downloadID != uuid.Nil {
			r.cache.Delete(buildKey(downloadID, userID, sessionID))
		}
		return nil
	}
	return r.Repository.Delete(fn)
}

// FindByDownloadID retrieves a mediaDownload by its downloadID from the repository.
func (r *MediaUserWatchPositionRepository) Find(
	downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID,
) (*ddownload.MediaUserWatchPosition, memsimple.CacheStatus, error) {
	find := func() (*ddownload.MediaUserWatchPosition, memsimple.CacheStatus, error) {
		position, status := r.cache.FindWithStatus(buildKey(downloadID, userID, sessionID))
		return position, status, nil
	}

	return r.Repository.FindWithStatus(find)
}

// Checks if a mediaDownload exists by its downloadID.
func (r *MediaUserWatchPositionRepository) Exists(
	downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID,
) (bool, error) {
	exists := func() (bool, error) {
		return r.cache.Exists(buildKey(downloadID, userID, sessionID)), nil
	}
	return r.Repository.Exists(exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *MediaUserWatchPositionRepository) CleanExpired(context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cache.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(clean)
}
