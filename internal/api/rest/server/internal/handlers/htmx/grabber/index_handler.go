package grabberh

import (
	"fmt"
	"time"

	htmxvalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/htmx/values"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

// IndexHandlers serves the main page (index.html)
func (h *GrabberHandlers) IndexHandler(ctx *fasthttp.RequestCtx) {
	var (
		needLoadHistory bool
	)

	resps, err := h.usecases.Downloader.LoadHistory(ctx, time.Now(), 1)
	if err == nil {
		needLoadHistory = len(resps) != 0
	}

	ctx.Response.Header.SetCookie(cookiePageHasDivItemsKey.makeCookie(fmt.Sprint(needLoadHistory), "/", 7*24*60*60))

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	dataMap := htmxvalues.MergeMaps(htmxvalues.IndexValues, htmxvalues.PathValues)
	dataMap[htmxvalues.NeedLoadHistoryKey] = needLoadHistory

	// Execute template with PageTitle
	if err := h.templates.ExecuteTemplate(ctx, htmxvalues.IndexHtmlFileName, dataMap); err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
		return
	}
}
