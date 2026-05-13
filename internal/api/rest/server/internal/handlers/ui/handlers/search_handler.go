package handlers

import (
	"bytes"
	"html/template"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) SearchHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)

	filters := make(requestFilters)

	filterByTitle := string(ctx.PostArgs().Peek(searchKey))
	if filterByTitle != "" {
		filters[filterByTitleKey] = filterByTitle
	}

	var rowsBuf bytes.Buffer
	err := h.getDownloadsHistory(ctx, &rowsBuf, *ctxUser, time.Now().UTC(), filters)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	extraData := make(map[string]any)
	extraData[uivalues.ResultNoRowsKey] = rowsBuf.Len() == 0
	extraData[uivalues.ResultRowsHTMLKey] = template.HTML(rowsBuf.String())

	pageData := uivalues.RowFragmentData{
		BasePaths: uivalues.NewBasePaths(),
		Values:    uivalues.RowFragmentValues{},
		Extra:     extraData,
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.ExecuteTemplate(&bodyBuffer, uivalues.ComponentResultRowsKey, pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}
