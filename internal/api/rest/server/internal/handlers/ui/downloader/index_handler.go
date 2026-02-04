package downloaderh

import (
	"fmt"
	"time"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

// IndexHandlers serves the main page (index.html)
func (h *DownloaderHandlers) IndexHandler(ctx *fasthttp.RequestCtx) {
	var (
		needLoadHistory bool
	)

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("authorization error: %v", err), fasthttp.StatusUnauthorized)
		return
	}

	resps, err := h.usecases.Downloader.LoadHistory(ctx, userID, time.Now(), 1)
	if err == nil {
		needLoadHistory = len(resps) != 0
	}

	ctx.Response.Header.SetCookie(cookiePageHasDivItemsKey.makeCookie(fmt.Sprint(needLoadHistory), "/", 7*24*60*60))

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	dataMap := uivalues.MergeMaps(uivalues.IndexValues, uivalues.FormGrabValues, uivalues.PathValues)
	dataMap[uivalues.NeedLoadHistoryKey] = needLoadHistory

	// Execute template with PageTitle
	if err := h.templates.ExecuteTemplate(ctx, uivalues.IndexHtmlFileName, dataMap); err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
		return
	}
}
