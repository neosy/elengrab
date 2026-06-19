package idto

import (
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
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
