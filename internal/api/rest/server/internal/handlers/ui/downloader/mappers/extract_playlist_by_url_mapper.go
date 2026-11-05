package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	ierrors "github.com/neosy/elengrab/internal/errors"
)

func (m *Mappers) MapPlaylistToExtractPlaylistByURLResponse(playlist *ucdto.Playlist) (dto.ExtractPlaylistByURLResponse, error) {
	if playlist == nil {
		return dto.ExtractPlaylistByURLResponse{}, ierrors.ErrFuncParamNullPointer
	}

	respItems := make([]dto.PlaylistItem, 0, len(playlist.Items))

	for _, item := range playlist.Items {
		respItems = append(
			respItems,
			dto.PlaylistItem{
				Index: item.Index,

				URL:   item.URL,
				Title: item.Title,

				DurationSec:  int(item.Duration.Seconds()),
				DurationText: item.DurationText(),
			},
		)
	}

	return dto.ExtractPlaylistByURLResponse{
		Items: respItems,
	}, nil
}
