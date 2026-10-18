package downloader

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemPageByDownloadIDHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		fasthttpx.WriteHTMLHeadResponse(ctx)
		return
	}

	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadID, err := h.extractDownloadID(ctx)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
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
