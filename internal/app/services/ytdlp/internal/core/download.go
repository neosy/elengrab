package core

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

func (c *Core) Download(
	ctx context.Context,
	url string,
	concurrentFragments uint8,
	options *dservices.DownloadOptions,
	downloadResultCh chan<- *ddownload.DownloadResult,
) {
	c.downloader.Download(
		ctx,
		url,
		concurrentFragments,
		options,
		c.getBestFormat,
		c.GetTitle,
		downloadResultCh,
	)
}
