package downloader

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (d *Downloader) FetchChannelInfo(
	ctx context.Context,
	mediaURL string,
	allowCookies bool,
) (*dtypes.ChannelSource, error) {
	var (
		cookieFilePath string
		err            error
	)

	if allowCookies && d.serviceOptions.AllowCookies {
		var err error

		cookieFilePath, err = helper.CookieFilePathFromURL(mediaURL, d.serviceOptions.CookiesDir)
		if err != nil {
			return nil, err
		}
	}

	var info *idto.ExtractInfo

	if cookieFilePath == "" {
		info, err = d.executor.FetchInfo(ctx, mediaURL)
		if err != nil {
			return nil, err
		}
	} else {
		info, err = d.executor.FetchInfo(ctx, mediaURL, idto.WithUseCookies(cookieFilePath))
		if err != nil {
			return nil, err
		}
	}

	if info == nil {
		return nil, nil
	}

	channel := d.buildChannel(info)

	return channel, nil
}

func (d *Downloader) FetchChannelInfoWithCookieFallback(
	ctx context.Context,
	mediaURL string,
) (*dtypes.ChannelSource, error) {
	info, err := d.executor.FetchInfo(ctx, mediaURL)
	if err != nil {
		if !d.serviceOptions.AllowCookies || !helper.CheckCookiesError(err) {
			return nil, err
		}

		cookieFilePath, err := helper.CookieFilePathFromURL(mediaURL, d.serviceOptions.CookiesDir)
		if err != nil {
			return nil, err
		}

		if cookieFilePath == "" {
			return nil, err
		}

		info, err = d.executor.FetchInfo(ctx, mediaURL, idto.WithUseCookies(cookieFilePath))
		if err != nil {
			return nil, err
		}
	}

	if info == nil {
		return nil, nil
	}

	channel := d.buildChannel(info)

	return channel, nil
}
