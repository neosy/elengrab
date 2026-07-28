package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
)

func (m *Mappers) MapMediaUserWatchChunkDomainToEntity(chunk *ddownload.MediaUserWatchChunk) (*ewatchevent.MediaUserWatchChunk, error) {
	return &ewatchevent.MediaUserWatchChunk{
		DownloadID: chunk.DownloadID,
		UserID:     chunk.UserID,
		ChunkIndex: int(chunk.ChunkIndex),
		Qty:        int(chunk.Qty),
	}, nil
}

func (m *Mappers) MapMediaUserWatchChunkEntityToDomain(chunk *ewatchevent.MediaUserWatchChunk) (*ddownload.MediaUserWatchChunk, error) {
	return &ddownload.MediaUserWatchChunk{
		DownloadID: chunk.DownloadID,
		UserID:     chunk.UserID,
		ChunkIndex: uint32(chunk.ChunkIndex),
		Qty:        uint32(chunk.Qty),
		CreatedAt:  chunk.CreatedAt,
	}, nil
}
