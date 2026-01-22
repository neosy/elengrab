package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/ports/persistence"
	nmemory "github.com/neosy/elengrab/pkg/ncache/memory"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type userIDFileIDKey struct {
	userID uuid.UUID
	fileID uuid.UUID
}

type DownloadStateRepository struct {
	nmemory.Repository[ddownload.DownloadState]

	userID *uuid.UUID

	dataByFileIdMap       nmemory.Cache[uuid.UUID, ddownload.DownloadState]
	dataByUserIDFileIDMap nmemory.Cache[userIDFileIDKey, ddownload.DownloadState]
	dataByTaskIdMap       nmemory.Cache[uuid.UUID, ddownload.DownloadState]
}

// newDownloadStateRepository returns a new object for the repository
func newDownloadStateRepository(ttl time.Duration) *DownloadStateRepository {
	r := &DownloadStateRepository{
		dataByFileIdMap:       make(nmemory.Cache[uuid.UUID, ddownload.DownloadState]),
		dataByUserIDFileIDMap: make(nmemory.Cache[userIDFileIDKey, ddownload.DownloadState]),
		dataByTaskIdMap:       make(nmemory.Cache[uuid.UUID, ddownload.DownloadState]),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *DownloadStateRepository) WithUser(userID uuid.UUID) persistence.DownloadStateRepository {
	return &DownloadStateRepository{
		Repository:            r.Repository,
		userID:                &userID,
		dataByFileIdMap:       r.dataByFileIdMap,
		dataByUserIDFileIDMap: r.dataByUserIDFileIDMap,
		dataByTaskIdMap:       r.dataByTaskIdMap,
	}
}

func (r *DownloadStateRepository) Save(_ context.Context, state *ddownload.DownloadState) error {
	if state == nil || state.FileId == uuid.Nil {
		return nil
	}

	save := func() error {
		if state.TaskId == nil {
			stateLast := r.dataByFileIdMap.Find(state.FileId, nil)
			if stateLast != nil && stateLast.TaskId != nil {
				r.dataByTaskIdMap.Delete(*stateLast.TaskId)
			}
		}

		stateCopy := r.copyDownloadState(state)

		r.dataByFileIdMap.Save(state.FileId, stateCopy, nil, r.TTL())

		var userID uuid.UUID
		if state.File != nil && state.File.UserID != nil {
			userID = *state.File.UserID
		}
		r.dataByUserIDFileIDMap.Save(userIDFileIDKey{userID, state.FileId}, stateCopy, nil, r.TTL())

		if state.TaskId != nil {
			r.dataByTaskIdMap.Save(*state.TaskId, stateCopy, nil, r.TTL())
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
		file := r.dataByFileIdMap.Find(fileId, nil)
		if file == nil {
			return nil
		}

		if file.TaskId != nil {
			r.dataByTaskIdMap.Delete(*file.TaskId)
		}

		var userID uuid.UUID
		if file.UserID != nil {
			userID = *file.UserID
		}
		r.dataByUserIDFileIDMap.Delete(userIDFileIDKey{userID, file.FileId})

		r.dataByFileIdMap.Delete(fileId)

		return nil
	}
	return r.Repository.Delete(delete)
}

func (r *DownloadStateRepository) FindByFileId(_ context.Context, fileId uuid.UUID) (*ddownload.DownloadState, error) {
	var find func() (*ddownload.DownloadState, error)

	if r.userID == nil {
		find = func() (*ddownload.DownloadState, error) {
			return r.dataByFileIdMap.Find(fileId, r.copyDownloadState), nil
		}
	} else {
		find = func() (*ddownload.DownloadState, error) {
			return r.dataByUserIDFileIDMap.Find(userIDFileIDKey{*r.userID, fileId}, r.copyDownloadState), nil
		}
	}
	return r.Repository.Find(find)
}

func (r *DownloadStateRepository) FindByTaskId(_ context.Context, taskId uuid.UUID) (*ddownload.DownloadState, error) {
	find := func() (*ddownload.DownloadState, error) {
		return r.dataByTaskIdMap.Find(taskId, r.copyDownloadState), nil
	}
	return r.Repository.Find(find)
}

func (r *DownloadStateRepository) copyDownloadState(state *ddownload.DownloadState) *ddownload.DownloadState {
	if state == nil {
		return nil
	}

	copy := uptr.Copy(state)
	copy.TaskId = uptr.Copy(state.TaskId)
	copy.Progress = uptr.Copy(state.Progress)

	if state.File != nil {
		copy.File = uptr.Copy(state.File)
		copy.File.UserID = uptr.Copy(state.File.UserID)
		copy.File.YoutubeChannelID = uptr.Copy(state.File.YoutubeChannelID)
		copy.File.FileSize = uptr.Copy(state.File.FileSize)
		copy.File.PartialHash = uptr.Copy(state.File.PartialHash)
		copy.File.MediaInfo = uptr.Copy(state.File.MediaInfo)
		copy.File.ErrorMessage = uptr.Copy(state.File.ErrorMessage)
		copy.File.DeletedAt = uptr.Copy(state.File.DeletedAt)
		copy.File.DownloadTask = uptr.Copy(state.File.DownloadTask)
	}

	return copy
}

func (r *DownloadStateRepository) CleanExpired(_ context.Context) error {
	clean := func() error {
		r.dataByFileIdMap.CleanExpired()
		r.dataByTaskIdMap.CleanExpired()
		r.dataByUserIDFileIDMap.CleanExpired()
		return nil
	}
	return r.Repository.CleanExpired(clean)
}
