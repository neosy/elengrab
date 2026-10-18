package ytdlpsrv

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func newRequestOptions(opts ...pservices.DownloaderOption) pservices.DownloaderOptions {
	options := pservices.DownloaderOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

func internalRequestOptions(opts ...pservices.DownloaderOption) idto.RequestOptions {
	options := newRequestOptions(opts...)
	iOptions := idto.DefaultRequestOptions()

	if options.AllowCookies() != nil {
		iOptions.AllowCookies = *options.AllowCookies()
	}

	return iOptions
}

func internalFetchImageOptions(opts ...pservices.DownloaderOption) idto.RequestOptions {
	options := newRequestOptions(opts...)
	iOptions := idto.DefaultFetchImageOptions()

	if options.AllowCookies() != nil {
		iOptions.AllowCookies = *options.AllowCookies()
	}

	if options.LimitBytes() != nil {
		iOptions.LimitBytes = *options.LimitBytes()
	}

	if options.Timeout() != nil {
		iOptions.Timeout = *options.Timeout()
	}

	return iOptions
}
