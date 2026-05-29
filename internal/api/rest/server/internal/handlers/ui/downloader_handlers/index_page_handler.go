package handlers

import (
	"mime"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/policy"
	"github.com/valyala/fasthttp"
)

// IndexHandlers serves the main page (index.html)
func (h *DownloaderHandlers) IndexPageHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		ctx.SetContentType(mime.TypeByExtension(".html"))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	if h.redirectGuestIfAuthRequired(ctx) {
		return
	}

	ctxUser := policy.ResolveUserOrAnonym(ctx)

	h.renderIndexPage(ctx, ctxUser)
}
