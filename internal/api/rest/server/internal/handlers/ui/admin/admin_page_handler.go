package admin

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

// AdminPageHandler serves the main page
func (h *AdminHandlers) AdminPageHandler(ctx *fasthttp.RequestCtx) {
	if err := fasthttpx.EnforceHTTPS(ctx); err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	if h.redirectAuth(ctx) {
		return
	}

	authCtx, err := policy.EnsureUser(ctx)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	if ctx.IsHead() {
		fasthttpx.WriteHTMLHeadResponse(ctx)
		return
	}

	h.renderPage(ctx, authCtx)
}
