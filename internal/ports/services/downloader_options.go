package pservices

import (
	"time"
)

type DownloaderOption func(opt *DownloaderOptions)

type DownloaderOptions struct {
	// AllowCookies allow cookies for YouTube, Instagram, etc
	allowCookies *bool

	limitBytes *int64
	timeout    *time.Duration
}

func (o *DownloaderOptions) AllowCookies() *bool {
	return o.allowCookies
}

func (o *DownloaderOptions) LimitBytes() *int64 {
	return o.limitBytes
}

func (o *DownloaderOptions) Timeout() *time.Duration {
	return o.timeout
}

// WithDownloaderCookies enables cookies for the request.
func WithDownloaderCookies() DownloaderOption {
	return func(opt *DownloaderOptions) {
		opt.allowCookies = new(true)
	}
}

// WithDownloaderLimitBytes sets the maximum number of bytes to fetch.
func WithDownloaderLimitBytes(bytes int64) DownloaderOption {
	return func(opt *DownloaderOptions) {
		opt.limitBytes = &bytes
	}
}

// WithDownloaderTimeout sets the maximum request execution duration.
func WithDownloaderTimeout(timeout time.Duration) DownloaderOption {
	return func(opt *DownloaderOptions) {
		opt.timeout = &timeout
	}
}
