package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) SearchHandler(ctx *fasthttp.RequestCtx) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetBodyString(fmt.Sprintf("Authorization error: %v", err))
		return
	}

	filters := make(requestFilters)

	filterByTitle := string(ctx.PostArgs().Peek(searchKey))
	if filterByTitle != "" {
		filters[filterByTitleKey] = filterByTitle
	}

	var rowsBuf bytes.Buffer
	err = h.getFilesHistory(ctx, &rowsBuf, userID, time.Now().UTC(), filters)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	dataMap := uivalues.MergeMaps(uivalues.PathValues)
	dataMap[uivalues.ResultNoRowsKey] = rowsBuf.Len() == 0
	dataMap[uivalues.ResultRowsHTMLKey] = template.HTML(rowsBuf.String())

	var bodyBuffer bytes.Buffer
	if err := h.templates.ExecuteTemplate(&bodyBuffer, uivalues.ResultRowsHtmlFileName, dataMap); err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}
