package downloader

import (
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ResolveShortLinkHandler(ctx *fasthttp.RequestCtx) {
	shortLink, downloadID, err := h.resolveShortLinkToDownloadID(ctx, true)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	h.renderMediaItemPage(ctx,
		renderWatchPageRequest{
			pageURL:    shortLink.shortURL,
			downloadID: downloadID,

			showBackButton: false,

			shortLink: shortLink,
		},
	)
}
