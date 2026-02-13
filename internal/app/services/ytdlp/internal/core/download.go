package core

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

func (c *Core) Download(
	ctx context.Context,
	url string,
	options *dservices.DownloadOptions,
	downloadResultCh chan<- *ddownload.DownloadResult,
) {
	c.downloader.Download(
		ctx,
		url,
		options,
		c.updateFormatCache,
		c.getBestFormat,
		c.GetTitle,
		downloadResultCh,
	)
}
