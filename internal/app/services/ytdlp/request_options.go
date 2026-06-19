package ytdlpsrv

import (
	"time"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

type RequestOption func(opt *RequestOptions)

type RequestOptions struct {
	// AllowCookies allow cookies for YouTube, Instagram, etc
	allowCookies *bool

	limitBytes *int64
	timeout    *time.Duration
}

func newRequestOptions(opts ...RequestOption) RequestOptions {
	options := RequestOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

func internalRequestOptions(opts ...RequestOption) idto.RequestOptions {
	options := newRequestOptions(opts...)
	iOptions := idto.DefaultRequestOptions()

	if options.allowCookies != nil {
		iOptions.AllowCookies = *options.allowCookies
	}

	return iOptions
}

func internalFetchImageOptions(opts ...RequestOption) idto.RequestOptions {
	options := newRequestOptions(opts...)
	iOptions := idto.DefaultFetchImageOptions()

	if options.allowCookies != nil {
		iOptions.AllowCookies = *options.allowCookies
	}

	if options.limitBytes != nil {
		iOptions.LimitBytes = *options.limitBytes
	}

	if options.timeout != nil {
		iOptions.Timeout = *options.timeout
	}

	return iOptions
}

// WithRequestCookies enables cookies for the request.
func WithRequestCookies() RequestOption {
	return func(opt *RequestOptions) {
		opt.allowCookies = new(true)
	}
}

// WithRequestLimitBytes sets the maximum number of bytes to fetch.
func WithRequestLimitBytes(bytes int64) RequestOption {
	return func(opt *RequestOptions) {
		opt.limitBytes = &bytes
	}
}

// WithRequestTimeout sets the maximum request execution duration.
func WithRequestTimeout(timeout time.Duration) RequestOption {
	return func(opt *RequestOptions) {
		opt.timeout = &timeout
	}
}
