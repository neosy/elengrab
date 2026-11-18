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

type uploadStates struct {
	items map[uuid.UUID]*ddownload.UploadState
}

func newUploadStates() *uploadStates {
	return &uploadStates{
		items: make(map[uuid.UUID]*ddownload.UploadState),
	}
}

func (c *uploadStates) Clone() *uploadStates {
	if c == nil {
		return nil
	}

	itemsCloned := maps.Clone(c.items)

	for k, v := range itemsCloned {
		itemsCloned[k] = v.Clone()
	}

	return &uploadStates{
		items: itemsCloned,
	}
}

type uploadStateRepositoryFilters struct {
	userID *uuid.UUID
}

type UploadStateRepository struct {
	memsimple.Repository[ddownload.UploadState]

	filters uploadStateRepositoryFilters

	cacheByDownloadID       memsimple.Cache[uuid.UUID, ddownload.UploadState]
	cacheFailedByDownloadID memsimple.Cache[uuid.UUID, ddownload.UploadState]

	cacheByUserID       memsimple.Cache[uuid.UUID, uploadStates]
	cacheFailedByUserID memsimple.Cache[uuid.UUID, uploadStates]
}

// newUploadStateRepository returns a new object for the repository
func newUploadStateRepository() *UploadStateRepository {
	r := &UploadStateRepository{
		cacheByDownloadID:       memsimple.NewCacheWithDeaultCloner[uuid.UUID, ddownload.UploadState, *ddownload.UploadState](),
		cacheFailedByDownloadID: memsimple.NewCacheWithDeaultCloner[uuid.UUID, ddownload.UploadState, *ddownload.UploadState](),

		cacheByUserID:       memsimple.NewCacheWithDeaultCloner[uuid.UUID, uploadStates, *uploadStates](),
		cacheFailedByUserID: memsimple.NewCacheWithDeaultCloner[uuid.UUID, uploadStates, *uploadStates](),
	}
	r.Repository.Init(0)
	return r
}

func (r *UploadStateRepository) Name() string {
	return "download_state"
}

func (r *UploadStateRepository) copy() *UploadStateRepository {
	return &UploadStateRepository{
		Repository:              r.Repository,
		cacheByDownloadID:       r.cacheByDownloadID,
		cacheFailedByDownloadID: r.cacheFailedByDownloadID,

		filters: r.filters,

		cacheByUserID:       r.cacheByUserID,
		cacheFailedByUserID: r.cacheFailedByUserID,
	}
}

func (r *UploadStateRepository) WithUser(userID uuid.UUID) persistence.UploadStateCacheRepository {
	copyRepo := r.copy()

	copyRepo.filters.userID = &userID

	return copyRepo
}

func (r *UploadStateRepository) saveToCacheByUser(
	state *ddownload.UploadState,
	cache *memsimple.Cache[uuid.UUID, uploadStates],
) {
	if state.IsEmpty() || cache == nil ||
		!state.Download.HasRegisteredUser() {
		return
	}

	userID := *state.Download.UserID

	userStates := cache.FindNoClone(userID)
	if userStates == nil {
		userStates = newUploadStates()
	}

	userStates.items[state.DownloadID()] = state

	cache.Save(userID, userStates, r.TTL())
}

func (r *UploadStateRepository) deleteFromCacheByUser(
	state *ddownload.UploadState,
	cache *memsimple.Cache[uuid.UUID, uploadStates],
) {
	if state.IsEmpty() || cache == nil ||
		!state.Download.HasRegisteredUser() {
		return
	}

	userID := state.UserID()

	userStates := cache.FindNoClone(userID)
	if userStates == nil {
		return
	}

	delete(userStates.items, state.DownloadID())

	if len(userStates.items) == 0 {
		cache.Delete(userID)
	}
}

func (r *UploadStateRepository) Save(ctx context.Context, state *ddownload.UploadState) error {
	if state.IsEmpty() {
		return nil
	}

	save := func() error {
		r.cacheByDownloadID.Save(state.DownloadID(), state, r.TTL())
		r.saveToCacheByUser(state, &r.cacheByUserID)

		if state.Status() == dtypes.MediaDownloadStatusFailed {
			r.cacheFailedByDownloadID.Save(state.DownloadID(), state, r.TTL())
			r.saveToCacheByUser(state, &r.cacheFailedByUserID)
		} else {
			r.cacheFailedByDownloadID.Delete(state.DownloadID())
			r.deleteFromCacheByUser(state, &r.cacheFailedByUserID)
		}

		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *UploadStateRepository) Delete(ctx context.Context, downloadID uuid.UUID) error {
	delete := func() error {
		state := r.cacheByDownloadID.Find(downloadID)
		if state == nil {
			return nil
		}

		r.deleteFromCacheByUser(state, &r.cacheByUserID)
		r.deleteFromCacheByUser(state, &r.cacheFailedByUserID)

		r.cacheByDownloadID.Delete(downloadID)
		r.cacheFailedByDownloadID.Delete(downloadID)

		return nil
	}
	return r.Repository.Delete(ctx, delete)
}

func (r *UploadStateRepository) Count(
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

func (r *UploadStateRepository) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.UploadState, error) {
	var find func() (*ddownload.UploadState, error)

	if r.filters.userID == nil {
		find = func() (*ddownload.UploadState, error) {
			return r.cacheByDownloadID.Find(downloadID), nil
		}
	} else {
		find = func() (*ddownload.UploadState, error) {
			userStates := r.cacheByUserID.Find(*r.filters.userID)
			if userStates == nil {
				return nil, nil
			}
			return userStates.items[downloadID], nil
		}
	}

	return r.Repository.Find(ctx, find)
}

// Clear clears all entries from the cache.
func (r *UploadStateRepository) Clear(ctx context.Context) error {
	clear := func() error {
		r.cacheByDownloadID.Clear()
		r.cacheFailedByDownloadID.Clear()

		r.cacheByUserID.Clear()
		r.cacheFailedByUserID.Clear()

		return nil
	}

	return r.Repository.Clear(ctx, clear)
}
