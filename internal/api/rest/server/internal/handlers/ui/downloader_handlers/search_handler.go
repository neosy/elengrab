package handlers

import (
	"bytes"
	"html/template"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/policy"
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
	extraData[items.ResultNoRowsKey] = rowsBuf.Len() == 0
	extraData[items.ResultRowsHTMLKey] = template.HTML(rowsBuf.String())

	pageData := pages.RowFragmentData{
		BasePaths: paths.NewPaths(),
		Values:    &pages.RowFragmentValues{},
		Extra:     extraData,
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.ExecuteTemplate(&bodyBuffer, components.ResultRowsKey, pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}
