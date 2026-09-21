package downloader

import (
	"fmt"
	"mime"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/pkg/debugx"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
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

	searchParameters, err := parseGetSearchParameters(ctx)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	paramValues, err := searchParameters.ParseValues()
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	query, err := h.mappers.MapSearchParameterValuesToUsecaseQuery(&paramValues)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	h.renderIndexPage(ctx, ctxUser, query)
}
