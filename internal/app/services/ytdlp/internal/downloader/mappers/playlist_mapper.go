package mappers

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	ierrors "github.com/neosy/elengrab/internal/errors"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (m *Mappers) MapPlaylistDtoToService(playlistDto *idto.Playlist) (*pservices.DownloaderPlaylist, error) {
	if playlistDto == nil {
		return nil, ierrors.ErrFuncParamNullPointer
	}

	items := make([]pservices.DownloaderPlaylistItem, 0, len(playlistDto.Items))

	for _, item := range playlistDto.Items {
		items = append(
			items,
			pservices.DownloaderPlaylistItem{
				ID:    item.ID,
				Index: item.Index,

				URL:   item.URL,
				Title: item.Title,

				Duration: item.Duration,
			},
		)
	}

	return &pservices.DownloaderPlaylist{
		Items: items,
	}, nil
}
