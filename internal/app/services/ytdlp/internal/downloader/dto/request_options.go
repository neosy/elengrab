package idto

import (
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/internal/app/utils/siteimage/channels"
	"github.com/neosy/elengrab/internal/app/utils/siteimage/thumbnails"
)

type RequestOptions struct {
	// AllowCookies allow cookies for YouTube, Instagram, etc
	AllowCookies bool

	LimitBytes int64
	Timeout    time.Duration
}

func DefaultRequestOptions() RequestOptions {
	return RequestOptions{
		AllowCookies: false,
	}
}

func DefaultFetchImageOptions() RequestOptions {
	return RequestOptions{
		AllowCookies: false,
		LimitBytes:   consts.FetchImageLimit,
		Timeout:      consts.FetchImageTimeout,
	}
}

func (o RequestOptions) ThumbnailFetchOption() thumbnails.FetchOption {
	options := thumbnails.FetchOptions{
		LimitBytes: o.LimitBytes,
		Timeout:    o.Timeout,
	}

	return thumbnails.FetchOptionWithOptions(options)
}

func (o RequestOptions) ChannelFetchOption() channels.FetchOption {
	options := thumbnails.FetchOptions{
		LimitBytes: o.LimitBytes,
		Timeout:    o.Timeout,
	}

	return channels.FetchOptionWithOptions(options)
}
