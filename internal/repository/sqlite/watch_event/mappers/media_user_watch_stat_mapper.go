package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
)

func (m *Mappers) MapMediaUserWatchStatDomainToEntity(stat *ddownload.MediaUserWatchStat) (*ewatchevent.MediaUserWatchStat, error) {
	return &ewatchevent.MediaUserWatchStat{
		DownloadID: stat.DownloadID,
		UserID:     stat.UserID,
		Views:      int(stat.Views),
	}, nil
}

func (m *Mappers) MapMediaUserWatchStatEntityToDomain(stat *ewatchevent.MediaUserWatchStat) (*ddownload.MediaUserWatchStat, error) {
	return &ddownload.MediaUserWatchStat{
		DownloadID: stat.DownloadID,
		UserID:     stat.UserID,
		Views:      uint32(stat.Views),
		UpdatedAt:  stat.UpdatedAt,
	}, nil
}
