package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/ports/persistence"
	nmemory "github.com/neosy/elengrab/pkg/ncache/memory"
)

type userIDFileIDKey struct {
	userID uuid.UUID
	fileID uuid.UUID
}

type DownloadStateRepository struct {
	nmemory.Repository[ddownload.DownloadState]

	userID *uuid.UUID

	cacheByFileId       nmemory.Cache[uuid.UUID, ddownload.DownloadState]
	cacheByUserIDFileID nmemory.Cache[userIDFileIDKey, ddownload.DownloadState]
	cacheByTaskId       nmemory.Cache[uuid.UUID, ddownload.DownloadState]
}

// newDownloadStateRepository returns a new object for the repository
func newDownloadStateRepository(ttl time.Duration) *DownloadStateRepository {
	r := &DownloadStateRepository{
		cacheByFileId:       nmemory.NewCacheWithDeaultCopier[uuid.UUID, ddownload.DownloadState, *ddownload.DownloadState](),
		cacheByUserIDFileID: nmemory.NewCacheWithDeaultCopier[userIDFileIDKey, ddownload.DownloadState, *ddownload.DownloadState](),
		cacheByTaskId:       nmemory.NewCacheWithDeaultCopier[uuid.UUID, ddownload.DownloadState, *ddownload.DownloadState](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *DownloadStateRepository) WithUser(userID uuid.UUID) persistence.DownloadStateRepository {
	return &DownloadStateRepository{
		Repository:          r.Repository,
		userID:              &userID,
		cacheByFileId:       r.cacheByFileId,
		cacheByUserIDFileID: r.cacheByUserIDFileID,
		cacheByTaskId:       r.cacheByTaskId,
	}
}

func (r *DownloadStateRepository) Save(_ context.Context, state *ddownload.DownloadState) error {
	if state == nil || state.FileId == uuid.Nil {
		return nil
	}

	save := func() error {
		if state.TaskId == nil {
			stateLast := r.cacheByFileId.Find(state.FileId)
			if stateLast != nil && stateLast.TaskId != nil {
				r.cacheByTaskId.Delete(*stateLast.TaskId)
			}
		}

		r.cacheByFileId.Save(state.FileId, state, r.TTL())

		var userID uuid.UUID
		if state.File != nil && state.File.UserID != nil {
			userID = *state.File.UserID
		}
		r.cacheByUserIDFileID.Save(userIDFileIDKey{userID, state.FileId}, state, r.TTL())

		if state.TaskId != nil {
			r.cacheByTaskId.Save(*state.TaskId, state, r.TTL())
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

func (r *DownloadStateRepository) Delete(_ context.Context, fileId uuid.UUID) error {
	delete := func() error {
		file := r.cacheByFileId.Find(fileId)
		if file == nil {
			return nil
		}

		if file.TaskId != nil {
			r.cacheByTaskId.Delete(*file.TaskId)
		}

		var userID uuid.UUID
		if file.UserID != nil {
			userID = *file.UserID
		}
		r.cacheByUserIDFileID.Delete(userIDFileIDKey{userID, file.FileId})

		r.cacheByFileId.Delete(fileId)

		return nil
	}
	return r.Repository.Delete(delete)
}

func (r *DownloadStateRepository) FindByFileId(_ context.Context, fileId uuid.UUID) (*ddownload.DownloadState, error) {
	var find func() (*ddownload.DownloadState, error)

	if r.userID == nil {
		find = func() (*ddownload.DownloadState, error) {
			return r.cacheByFileId.Find(fileId), nil
		}
	} else {
		find = func() (*ddownload.DownloadState, error) {
			return r.cacheByUserIDFileID.Find(userIDFileIDKey{*r.userID, fileId}), nil
		}
	}
	return r.Repository.Find(find)
}

func (r *DownloadStateRepository) FindByTaskId(_ context.Context, taskId uuid.UUID) (*ddownload.DownloadState, error) {
	find := func() (*ddownload.DownloadState, error) {
		return r.cacheByTaskId.Find(taskId), nil
	}
	return r.Repository.Find(find)
}

func (r *DownloadStateRepository) CleanExpired(_ context.Context) error {
	clean := func() error {
		r.cacheByFileId.CleanExpired()
		r.cacheByTaskId.CleanExpired()
		r.cacheByUserIDFileID.CleanExpired()
		return nil
	}
	return r.Repository.CleanExpired(clean)
}
