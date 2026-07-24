package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type userIDDownloadIDKey struct {
	userID     uuid.UUID
	downloadID uuid.UUID
}

type DownloadStateRepository struct {
	memsimple.Repository[ddownload.DownloadState]

	userID *uuid.UUID

	cacheByDownloadID       memsimple.Cache[uuid.UUID, ddownload.DownloadState]
	cacheByUserIDDownloadID memsimple.Cache[userIDDownloadIDKey, ddownload.DownloadState]
	cacheByTaskID           memsimple.Cache[uuid.UUID, ddownload.DownloadState]
}

// newDownloadStateRepository returns a new object for the repository
func newDownloadStateRepository(ttl time.Duration) *DownloadStateRepository {
	r := &DownloadStateRepository{
		cacheByDownloadID:       memsimple.NewCacheWithDeaultCopier[uuid.UUID, ddownload.DownloadState, *ddownload.DownloadState](),
		cacheByUserIDDownloadID: memsimple.NewCacheWithDeaultCopier[userIDDownloadIDKey, ddownload.DownloadState, *ddownload.DownloadState](),
		cacheByTaskID:           memsimple.NewCacheWithDeaultCopier[uuid.UUID, ddownload.DownloadState, *ddownload.DownloadState](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *DownloadStateRepository) WithUser(userID uuid.UUID) persistence.DownloadStateCacheRepository {
	return &DownloadStateRepository{
		Repository:              r.Repository,
		userID:                  &userID,
		cacheByDownloadID:       r.cacheByDownloadID,
		cacheByUserIDDownloadID: r.cacheByUserIDDownloadID,
		cacheByTaskID:           r.cacheByTaskID,
	}
}

func (r *DownloadStateRepository) Save(_ context.Context, state *ddownload.DownloadState) error {
	if state == nil || state.DownloadID == uuid.Nil {
		return nil
	}

	save := func() error {
		if state.TaskID == nil {
			stateLast := r.cacheByDownloadID.Find(state.DownloadID)
			if stateLast != nil && stateLast.TaskID != nil {
				r.cacheByTaskID.Delete(*stateLast.TaskID)
			}
		}

		r.cacheByDownloadID.Save(state.DownloadID, state, r.TTL())

		var userID uuid.UUID
		if state.Download != nil && state.Download.UserID != nil {
			userID = *state.Download.UserID
		}
		r.cacheByUserIDDownloadID.Save(userIDDownloadIDKey{userID, state.DownloadID}, state, r.TTL())

		if state.TaskID != nil {
			r.cacheByTaskID.Save(*state.TaskID, state, r.TTL())
		}

		return nil
	}

	return r.Repository.Save(save)
}

func (r *DownloadStateRepository) Insert(ctx context.Context, state *ddownload.DownloadState) error {
	return r.Save(ctx, state)
}

func (r *DownloadStateRepository) Update(ctx context.Context, state *ddownload.DownloadState) error {
	return r.Save(ctx, state)
}

func (r *DownloadStateRepository) Delete(_ context.Context, downloadID uuid.UUID) error {
	delete := func() error {
		download := r.cacheByDownloadID.Find(downloadID)
		if download == nil {
			return nil
		}

		if download.TaskID != nil {
			r.cacheByTaskID.Delete(*download.TaskID)
		}

		var userID uuid.UUID
		if download.UserID != nil {
			userID = *download.UserID
		}
		r.cacheByUserIDDownloadID.Delete(userIDDownloadIDKey{userID, download.DownloadID})

		r.cacheByDownloadID.Delete(downloadID)

		return nil
	}
	return r.Repository.Delete(delete)
}

func (r *DownloadStateRepository) FindByDownloadID(_ context.Context, downloadID uuid.UUID) (*ddownload.DownloadState, error) {
	var find func() (*ddownload.DownloadState, error)

	if r.userID == nil {
		find = func() (*ddownload.DownloadState, error) {
			return r.cacheByDownloadID.Find(downloadID), nil
		}
	} else {
		find = func() (*ddownload.DownloadState, error) {
			return r.cacheByUserIDDownloadID.Find(userIDDownloadIDKey{*r.userID, downloadID}), nil
		}
	}
	return r.Repository.Find(find)
}

func (r *DownloadStateRepository) FindByTaskID(_ context.Context, taskId uuid.UUID) (*ddownload.DownloadState, error) {
	find := func() (*ddownload.DownloadState, error) {
		return r.cacheByTaskID.Find(taskId), nil
	}
	return r.Repository.Find(find)
}

func (r *DownloadStateRepository) CleanExpired(_ context.Context) error {
	clean := func() error {
		r.cacheByDownloadID.CleanExpired()
		r.cacheByTaskID.CleanExpired()
		r.cacheByUserIDDownloadID.CleanExpired()
		return nil
	}
	return r.Repository.CleanExpired(clean)
}
