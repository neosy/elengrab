package downloader

import (
	"bytes"
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/consts"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) SearchHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	filters := make(dtypes.QueryFiltersByName)

	var viewMode = dtypes.QueryMediaViewModeDefault

	viewModeStr := string(ctx.PostArgs().Peek(qkeys.ViewModeKey.String()))
	if viewModeStr != "" {
		var err error
		viewMode, err = dtypes.ParseQueryMediaViewMode(viewModeStr)
		if err != nil {
			fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}
	}

	searchText := types.SearchText(ctx.PostArgs().Peek(qkeys.SearchKey.String()))
	if searchText.IsLongEnough() {
		if err := searchText.Validate(); err != nil {
			fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}

		filters.Add(dtypes.QueryFilterNameSearch, searchText.String())
	}

	query := udto.MediaDownloadQueryDefault(consts.LoadHistoryLimit)
	query.ViewMode = viewMode
	query.Filters = filters

	var rowsBuf bytes.Buffer
	err := h.listDownloadsItems(ctx, &rowsBuf, authCtx, query)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
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

	var bodyBuffer bytes.Buffer
	if err := h.templates.Base.ExecuteTemplate(&bodyBuffer, components.ResultRowsKey, pageData); err != nil {
		fasthttpx.WriteErrorx(ctx, errInternal(err))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}
