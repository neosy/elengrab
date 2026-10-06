package downloader

import (
	"mime"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemPageByDownloadIDHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		ctx.SetContentType(mime.TypeByExtension(".html"))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadID, err := h.extractDownloadID(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	h.renderMediaItemPage(ctx,
		renderWatchPageRequest{
			pageURL:        httppaths.BuildMediaItemPath(downloadID),
			downloadID:     downloadID,
			showBackButton: true,
			authCtx:        authCtx,
		},
	)
}
