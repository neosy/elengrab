package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (uc *downloader) ExtractPlaylistByURL(
	ctx context.Context,
	authCtx dauth.AuthContext,
	url string,
) (*dto.Playlist, error) {
	playListFetched, err := uc.downloaderSrv.FetchPlaylist(ctx, url)
	if err != nil {
		return nil, err
	}

	playlist, err := uc.mappers.MapPlaylisServiceToDto(playListFetched)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}
