package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
)

func (m *Mappers) MapMediaWatchStatDomainToEntity(stat *ddownload.MediaWatchStat) (*ewatchevent.MediaWatchStat, error) {
	return &ewatchevent.MediaWatchStat{
		DownloadID: stat.DownloadID,
		Views:      int(stat.Views),
	}, nil
}

func (m *Mappers) MapMediaWatchStatEntityToDomain(stat *ewatchevent.MediaWatchStat) (*ddownload.MediaWatchStat, error) {
	return &ddownload.MediaWatchStat{
		DownloadID: stat.DownloadID,
		Views:      uint32(stat.Views),
		UpdatedAt:  stat.UpdatedAt,
	}, nil
}
