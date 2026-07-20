package downloader

import (
	"fmt"
	"mime"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/pkg/debugx"
	"github.com/valyala/fasthttp"
)

// IndexHandlers serves the main page (index.html)
func (h *DownloaderHandlers) IndexPageHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		ctx.SetContentType(mime.TypeByExtension(".html"))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	defer debugx.DumpGoroutinesIfTimeout(
		fmt.Sprintf("IndexPageHandler [%s %s]", ctx.Method(), ctx.Path()),
		5*time.Second,
	)()

	if h.redirectGuestIfAuthRequired(ctx) {
		return
	}

	ctxUser := policy.ResolveUserOrAnonym(ctx)

	h.renderIndexPage(ctx, ctxUser)
}
