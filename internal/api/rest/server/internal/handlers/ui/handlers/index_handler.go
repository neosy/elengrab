package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

// IndexHandlers serves the main page (index.html)
func (h *DownloaderHandlers) IndexHandler(ctx *fasthttp.RequestCtx) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("authorization error: %v", err), fasthttp.StatusUnauthorized)
		return
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	var rowsBuf bytes.Buffer
	err = h.getFilesHistory(ctx, &rowsBuf, userID, time.Now().UTC(), nil)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	showHistorySearch := h.usecases.Downloader.HistoryMode() != dtypes.HistoryModeDisabled

	dataMap := uivalues.MergeMaps(uivalues.IndexValues, uivalues.FormGrabValues, uivalues.PathValues)
	dataMap[uivalues.ShowHistorySearchKey] = showHistorySearch
	dataMap[uivalues.ResultNoRowsKey] = rowsBuf.Len() == 0
	dataMap[uivalues.ResultRowsHTMLKey] = template.HTML(rowsBuf.String())

	// Execute template with PageTitle
	if err := h.templates.ExecuteTemplate(ctx, uivalues.IndexHtmlFileName, dataMap); err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
		return
	}
}
