package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (m *Mappers) MapPlaylisServiceToDto(playlistIn *pservices.DownloaderPlaylist) (*dto.Playlist, error) {
	if playlistIn.Count() == 0 {
		return nil, nil
	}

	itemsOut := make([]dto.PlaylistItem, 0, len(playlistIn.Items))

	for _, item := range playlistIn.Items {
		itemsOut = append(
			itemsOut,
			dto.PlaylistItem{
				URL:   item.URL,
				Title: item.Title,
			},
		)
	}

	return &dto.Playlist{
		Items: itemsOut,
	}, nil
}
