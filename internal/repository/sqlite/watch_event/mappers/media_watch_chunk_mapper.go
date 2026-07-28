package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
)

func (m *Mappers) MapMediaWatchChunkDomainToEntity(chunk *ddownload.MediaWatchChunk) (*ewatchevent.MediaWatchChunk, error) {
	return &ewatchevent.MediaWatchChunk{
		DownloadID: chunk.DownloadID,
		UserID:     chunk.UserID,
		ChunkIndex: int(chunk.ChunkIndex),
		Qty:        int(chunk.Qty),
	}, nil
}

func (m *Mappers) MapMediaWatchChunkEntityToDomain(chunk *ewatchevent.MediaWatchChunk) (*ddownload.MediaWatchChunk, error) {
	return &ddownload.MediaWatchChunk{
		DownloadID: chunk.DownloadID,
		UserID:     chunk.UserID,
		ChunkIndex: uint32(chunk.ChunkIndex),
		Qty:        uint32(chunk.Qty),
		CreatedAt:  chunk.CreatedAt,
	}, nil
}
