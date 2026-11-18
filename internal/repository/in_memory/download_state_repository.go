package inmemory

import (
	"context"
	"maps"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type downloadStates struct {
	items map[uuid.UUID]*ddownload.DownloadState
}

func newDownloadStates() *downloadStates {
	return &downloadStates{
		items: make(map[uuid.UUID]*ddownload.DownloadState),
	}
}

func (c *downloadStates) Clone() *downloadStates {
	if c == nil {
		return nil
	}

	itemsCloned := maps.Clone(c.items)

	for k, v := range itemsCloned {
		itemsCloned[k] = v.Clone()
	}

	return &downloadStates{
		items: itemsCloned,
	}
}

type downloadStateRepositoryFilters struct {
	userID *uuid.UUID
}

type DownloadStateRepository struct {
	memsimple.Repository[ddownload.DownloadState]

	filters uploadStateRepositoryFilters

	cacheByDownloadID       memsimple.Cache[uuid.UUID, ddownload.DownloadState]
	cacheFailedByDownloadID memsimple.Cache[uuid.UUID, ddownload.DownloadState]
	cacheByTaskID           memsimple.Cache[uuid.UUID, ddownload.DownloadState]

	cacheByUserID       memsimple.Cache[uuid.UUID, downloadStates]
	cacheFailedByUserID memsimple.Cache[uuid.UUID, downloadStates]
}

// newDownloadStateRepository returns a new object for the repository
func newDownloadStateRepository() *DownloadStateRepository {
	r := &DownloadStateRepository{
		cacheByDownloadID:       memsimple.NewCacheWithDeaultCloner[uuid.UUID, ddownload.DownloadState, *ddownload.DownloadState](),
		cacheFailedByDownloadID: memsimple.NewCacheWithDeaultCloner[uuid.UUID, ddownload.DownloadState, *ddownload.DownloadState](),
		cacheByTaskID:           memsimple.NewCacheWithDeaultCloner[uuid.UUID, ddownload.DownloadState, *ddownload.DownloadState](),

		cacheByUserID:       memsimple.NewCacheWithDeaultCloner[uuid.UUID, downloadStates, *downloadStates](),
		cacheFailedByUserID: memsimple.NewCacheWithDeaultCloner[uuid.UUID, downloadStates, *downloadStates](),
	}
	r.Repository.Init(0)
	return r
}

func (r *DownloadStateRepository) Name() string {
	return "download_state"
}

func (r *DownloadStateRepository) copy() *DownloadStateRepository {
	return &DownloadStateRepository{
		Repository: r.Repository,
		filters:    r.filters,

		cacheByDownloadID:       r.cacheByDownloadID,
		cacheFailedByDownloadID: r.cacheFailedByDownloadID,
		cacheByTaskID:           r.cacheByTaskID,

		cacheByUserID:       r.cacheByUserID,
		cacheFailedByUserID: r.cacheFailedByUserID,
	}
}

func (r *DownloadStateRepository) WithUser(userID uuid.UUID) persistence.DownloadStateCacheRepository {
	copyRepo := r.copy()

	copyRepo.filters.userID = &userID

	return copyRepo
}

func (r *DownloadStateRepository) saveToCacheByUser(
	state *ddownload.DownloadState,
	cache *memsimple.Cache[uuid.UUID, downloadStates],
) {
	if state.IsEmpty() || cache == nil ||
		!state.Download.HasRegisteredUser() {
		return
	}

	userID := state.Download.EnsureUserID()

	userStates := cache.FindNoClone(userID)
	if userStates == nil {
		userStates = newDownloadStates()
	}

	userStates.items[state.DownloadID()] = state

	cache.Save(userID, userStates, r.TTL())
}

func (r *DownloadStateRepository) deleteFromCacheByUser(
	state *ddownload.DownloadState,
	cache *memsimple.Cache[uuid.UUID, downloadStates],
) {
	if state.IsEmpty() || cache == nil ||
		!state.Download.HasRegisteredUser() {
		return
	}

	userID := *state.Download.UserID

	userStates := cache.FindNoClone(userID)
	if userStates == nil {
		return
	}

	delete(userStates.items, state.DownloadID())

	if len(userStates.items) == 0 {
		cache.Delete(userID)
	}
}

func (r *DownloadStateRepository) Save(ctx context.Context, state *ddownload.DownloadState) error {
	if state == nil || state.Download == nil {
		return nil
	}

	save := func() error {
		if state.TaskID == nil {
			stateLast := r.cacheByDownloadID.Find(state.DownloadID())
			if stateLast != nil && stateLast.TaskID != nil {
				r.cacheByTaskID.Delete(*stateLast.TaskID)
			}
		}

		r.cacheByDownloadID.Save(state.DownloadID(), state, r.TTL())
		r.saveToCacheByUser(state, &r.cacheByUserID)

		if state.DownloadStatus() == dtypes.MediaDownloadStatusFailed {
			r.cacheFailedByDownloadID.Save(state.DownloadID(), state, r.TTL())
			r.saveToCacheByUser(state, &r.cacheFailedByUserID)
		} else {
			r.cacheFailedByDownloadID.Delete(state.DownloadID())
			r.deleteFromCacheByUser(state, &r.cacheFailedByUserID)
		}

		if state.TaskID != nil {
			r.cacheByTaskID.Save(*state.TaskID, state, r.TTL())
		}

		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *DownloadStateRepository) Insert(ctx context.Context, state *ddownload.DownloadState) error {
	return r.Save(ctx, state)
}

func (r *DownloadStateRepository) Update(ctx context.Context, state *ddownload.DownloadState) error {
	return r.Save(ctx, state)
}

func (r *DownloadStateRepository) Delete(ctx context.Context, downloadID uuid.UUID) error {
	delete := func() error {
		state := r.cacheByDownloadID.Find(downloadID)
		if state == nil {
			return nil
		}

		if state.TaskID != nil {
			r.cacheByTaskID.Delete(*state.TaskID)
		}

		r.deleteFromCacheByUser(state, &r.cacheByUserID)
		r.deleteFromCacheByUser(state, &r.cacheFailedByUserID)

		r.cacheByDownloadID.Delete(downloadID)
		r.cacheFailedByDownloadID.Delete(downloadID)

		return nil
	}
	return r.Repository.Delete(ctx, delete)
}

func (r *DownloadStateRepository) Count(
	ctx context.Context,
	userID uuid.UUID,
) (dtypes.IncompleteDownloadsCount, error) {
	var (
		downloadsCount dtypes.IncompleteDownloadsCount
		err            error
	)

	_, err = r.Repository.Len(
		ctx,
		func() (int, error) {
			userStates := r.cacheByUserID.Find(userID)
			if userStates != nil {
				downloadsCount.UserCount = len(userStates.items)
			}

			userStates = r.cacheFailedByUserID.Find(userID)
			if userStates != nil {
				downloadsCount.UserFailedCount = len(userStates.items)
			}

			downloadsCount.UserInProgressCount = downloadsCount.UserCount - downloadsCount.UserFailedCount

			downloadsCount.TotalCount = r.cacheByDownloadID.Len()
			downloadsCount.TotalFailedCount = r.cacheFailedByDownloadID.Len()
			downloadsCount.TotalInProgressCount = downloadsCount.TotalCount - downloadsCount.TotalFailedCount

			return downloadsCount.TotalCount, nil
		},
	)

	return downloadsCount, err
}

func (r *DownloadStateRepository) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.DownloadState, error) {
	var find func() (*ddownload.DownloadState, error)

	if r.filters.userID == nil {
		find = func() (*ddownload.DownloadState, error) {
			return r.cacheByDownloadID.Find(downloadID), nil
		}
	} else {
		find = func() (*ddownload.DownloadState, error) {
			userStates := r.cacheByUserID.Find(*r.filters.userID)
			if userStates == nil {
				return nil, nil
			}
			return userStates.items[downloadID], nil
		}
	}

	return r.Repository.Find(ctx, find)
}

func (r *DownloadStateRepository) FindByTaskID(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadState, error) {
	find := func() (*ddownload.DownloadState, error) {
		return r.cacheByTaskID.Find(taskId), nil
	}
	return r.Repository.Find(ctx, find)
}

// Clear clears all entries from the cache.
func (r *DownloadStateRepository) Clear(ctx context.Context) error {
	clear := func() error {
		r.cacheByDownloadID.Clear()
		r.cacheFailedByDownloadID.Clear()
		r.cacheByTaskID.Clear()

		r.cacheByUserID.Clear()
		r.cacheFailedByUserID.Clear()

		return nil
	}

	return r.Repository.Clear(ctx, clear)
}
