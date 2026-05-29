package handlers

import (
	"mime"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/policy"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

// AdminPageHandler serves the main page
func (h *AdminHandlers) AdminPageHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		ctx.SetContentType(mime.TypeByExtension(".html"))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	if err := nfasthttp.EnforceHTTPS(ctx); err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if h.redirectAuth(ctx) {
		return
	}

	ctxUser, err := policy.EnsureUser(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	h.renderPage(ctx, ctxUser)
}
