package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dconsts "github.com/neosy/elengrab/internal/domain/consts"
)

func (m *Mappers) MapPlaylistToExtractPlaylistByURLResponse(playlist *ucdto.Playlist) (dto.ExtractPlaylistByURLResponse, error) {
	var emptyResponse = dto.ExtractPlaylistByURLResponse{
		Items: []dto.PlaylistItem{},
	}

	if playlist.Count() == 0 {
		return emptyResponse, nil
	}

	items := playlist.Items
	truncated := len(items) > dconsts.MediaDownloadListMaxLen

	if truncated {
		items = items[:dconsts.MediaDownloadListMaxLen]
	}

	resultItems := make([]dto.PlaylistItem, 0, len(items))

	for _, item := range items {
		resultItems = append(
			resultItems,
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
		Items:     resultItems,
		Count:     len(resultItems),
		Truncated: truncated,
	}, nil
}
