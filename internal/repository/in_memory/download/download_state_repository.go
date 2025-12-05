package imdownload

import (
	"context"
	"sync"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/repository/in_memory/download/mappers"
)

type DownloadStateRepository struct {
	mappers *mappers.Mappers
	mu      sync.Mutex

	dataByFileIdMap  map[uuid.UUID]*ddownload.DownloadState
	dataByStateIdMap map[uuid.UUID]*ddownload.DownloadState
}

// NewDownloadStateRepository returns a new object for the repository
func NewDownloadStateRepository() *DownloadStateRepository {
	return &DownloadStateRepository{
		mappers:          mappers.NewMappers(),
		dataByFileIdMap:  make(map[uuid.UUID]*ddownload.DownloadState),
		dataByStateIdMap: make(map[uuid.UUID]*ddownload.DownloadState),
	}
}

func (r *DownloadStateRepository) Save(_ context.Context, state *ddownload.DownloadState) error {
	if state == nil || state.FileId == uuid.Nil {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if state.TaskId != nil && *state.TaskId == uuid.Nil {
		state.TaskId = nil
		if state.File != nil {
			state.File.DownloadTask = nil
		}
	}

	if state.File != nil && state.File.DownloadTask == nil && state.TaskId != nil {
		delete(r.dataByStateIdMap, *state.TaskId)
		state.TaskId = nil
	}

	r.dataByFileIdMap[state.FileId] = state

	if state.TaskId != nil && *state.TaskId != uuid.Nil {
		r.dataByStateIdMap[*state.TaskId] = state
	}

	return nil
}

func (r *DownloadStateRepository) Insert(ctx context.Context, state *ddownload.DownloadState) error {
	return r.Save(ctx, state)
}

func (r *DownloadStateRepository) Update(ctx context.Context, state *ddownload.DownloadState) error {
	return r.Save(ctx, state)
}

func (r *DownloadStateRepository) Delete(_ context.Context, fileId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, exists := r.dataByFileIdMap[fileId]
	if !exists || file == nil {
		return nil
	}
	if file.TaskId != nil && *file.TaskId != uuid.Nil {
		delete(r.dataByStateIdMap, *file.TaskId)
	}
	delete(r.dataByFileIdMap, fileId)

	return nil
}

func (r *DownloadStateRepository) FindByFileId(_ context.Context, fileId uuid.UUID) (*ddownload.DownloadState, error) {
	return r.dataByFileIdMap[fileId], nil
}

func (r *DownloadStateRepository) FindByTaskId(_ context.Context, taskId uuid.UUID) (*ddownload.DownloadState, error) {
	return r.dataByStateIdMap[taskId], nil
}
