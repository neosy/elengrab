package downloader

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (d *Downloader) FetchPlaylistWithCookieFallback(
	ctx context.Context,
	playlistURL string,
) (*pservices.DownloaderPlaylist, error) {
	playListDto, err := d.executor.FetchPlaylist(ctx, playlistURL)
	if err != nil {
		if !d.serviceOptions.AllowCookies || !helper.CheckCookiesError(err) {
			return nil, err
		}

		cookieFilePath, err := d.cookieFilePath(playlistURL)
		if err != nil {
			return nil, err
		}

		if cookieFilePath == "" {
			return nil, err
		}

		playListDto, err = d.executor.FetchPlaylist(ctx, playlistURL, idto.WithUseCookies(cookieFilePath))
		if err != nil {
			return nil, err
		}
	}

	if playListDto.Count() == 0 {
		return nil, nil
	}

	playlist, err := d.mappers.MapPlaylistDtoToService(playListDto)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}
