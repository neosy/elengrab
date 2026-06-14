package downloader

import (
	"bytes"
	"html/template"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) SearchHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	filters := make(requestFilters)

	filterByTitle := string(ctx.PostArgs().Peek(searchKey))
	if filterByTitle != "" {
		filters[filterByTitleKey] = filterByTitle
	}

	var rowsBuf bytes.Buffer
	err := h.getDownloadsHistory(ctx, &rowsBuf, authCtx, time.Now().UTC(), filters)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	extraData := make(map[string]any)
	extraData[items.ResultNoRowsKey] = rowsBuf.Len() == 0
	extraData[items.ResultRowsHTMLKey] = template.HTML(rowsBuf.String())

	pageData := pages.RowFragmentData{
		BasePaths: paths.NewHttpPaths(),
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
