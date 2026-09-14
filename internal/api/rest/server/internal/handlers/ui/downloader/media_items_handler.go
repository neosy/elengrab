package downloader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/consts"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemsHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	if ctx.IsGet() {
		h.mediaItemsGet(ctx)
		return
	}

	if ctx.IsPost() {
		h.mediaItemsPost(ctx)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
}

func (h *DownloaderHandlers) mediaItemsGet(ctx *fasthttp.RequestCtx) {
	var (
		searchParameters *types.SearchParameters
	)

	ctxUser := policy.ResolveUserOrAnonym(ctx)

	searchParameters, err := parseGetSearchParameters(ctx)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	filters, err := parseGetFilters(ctx)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	viewMode := dtypes.QueryMediaViewModeDefault
	lastRecord := dtypes.QueryMediaDownloadCursor{}

	if searchParameters != nil {
		lastQueryCursor, err := searchParameters.ParseLastCursor()
		if err != nil {
			fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}

		if lastQueryCursor != nil {
			viewMode = lastQueryCursor.ViewMode
			lastRecord = dtypes.QueryMediaDownloadCursor{
				ID:        lastQueryCursor.LastID,
				CreatedAt: lastQueryCursor.LastCreateAt,
				Views:     lastQueryCursor.LastViews,
			}
		}
	}

	query := udto.BuildMediaDownloadQuery(
		udto.MediaDownloadQuery{
			ViewMode:   viewMode,
			Limit:      consts.LoadHistoryLimit,
			LastRecord: lastRecord,
			Filters:    filters,
		},
	)

	var bodyBuffer bytes.Buffer
	err = h.listDownloadsItems(ctx, &bodyBuffer, ctxUser, query)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}

func (h *DownloaderHandlers) mediaItemsPost(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)

	var postReq dto.MediaItemsRequest
	err := json.Unmarshal(ctx.PostBody(), &postReq)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	err = h.validators.Validate.Struct(postReq)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	query, err := h.mappers.MapMediaItemsRequestToQuery(postReq)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	var bodyBuffer bytes.Buffer
	err = h.listDownloadsItems(ctx, &bodyBuffer, ctxUser, query)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyString("")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}

func (h *DownloaderHandlers) listDownloadsItems(
	ctx context.Context,
	buf *bytes.Buffer,
	authCtx dauth.AuthContext,
	query udto.MediaDownloadQuery,
) error {
	query.Limit++

	downloads, err := h.downloader.ListDownloadInfo(ctx, authCtx, query)
	if err != nil {
		return err
	}

	if len(downloads) == 0 {
		return nil
	}

	shouldLoadNextHistory := len(downloads) > consts.LoadHistoryLimit

	// If there are more items than the limit, we show only the limited number of items and a "Load more"
	lines := downloads
	if len(downloads) > consts.LoadHistoryLimit {
		lines = downloads[:consts.LoadHistoryLimit]
	}

	lastLine := lines[len(lines)-1]
	lastDownloadID := lastLine.DownloadID
	lastCreatedAt := lastLine.CreatedAt
	lastViews := lastLine.ViewCount

	for i, downloadInfo := range lines {
		row := h.renderMediaItemRow(
			ctx,
			renderMediaItemRowParams{
				downloadInfo:   downloadInfo,
				lazyLoadImages: true,
			},
		)
		if row.err != nil {
			h.logger.Warn(
				"Failed to generate row",
				"error", err,
			)
			continue
		}

		err = h.templates.Base.ExecuteTemplate(buf, components.ResultRowStatusKey, row.data)
		if err != nil {
			h.logger.Warn(
				"Failed to execute template",
				"name", components.ResultRowStatusKey,
				"error", err,
			)
			continue
		}

		if shouldLoadNextHistory && i == consts.PreloadHistoryAfter-1 {
			query := udto.BuildMediaDownloadQuery(query)
			query.LastRecord.ID = lastDownloadID
			query.LastRecord.CreatedAt = lastCreatedAt
			query.LastRecord.Views = lastViews
			h.renderRowShouldLoadHistory(buf, query)
		}
	}

	if shouldLoadNextHistory {
		h.genRowLoadHistory(buf)
	}

	return nil
}

func (h *DownloaderHandlers) genRowLoadHistory(buf *bytes.Buffer) error {
	extraData := make(map[string]any)
	extraData[items.DisableHTMXEventKey] = true

	pageData := pages.PageFragmentData{
		Extra: extraData,
	}

	err := h.templates.Base.ExecuteTemplate(buf, components.ResultLoadHistoryKey, pageData)
	if err != nil {
		return errInternal(err)
	}

	return nil
}

func (h *DownloaderHandlers) renderRowShouldLoadHistory(
	buf *bytes.Buffer,
	query udto.MediaDownloadQuery,
) error {
	if query.LastRecord.CreatedAt.IsZero() {
		return nil
	}

	getSearchText := func(filters dtypes.QueryFiltersByName) dtypes.SearchText {
		if len(filters) == 0 {
			return ""
		}

		filter, exists := filters[dtypes.QueryFilterNameSearch]
		if !exists {
			return ""
		}

		txt, ok := filter.Condition().Value().(string)
		if !ok {
			return ""
		}

		return dtypes.SearchText(txt)
	}

	var queryParameters []string

	if text := getSearchText(query.Filters); text.IsValidate() {
		queryParameters = append(queryParameters, fmt.Sprintf("filter[%s]=%s", qkeys.SearchKey.String(), text.String()))
	}

	lastCursor := types.MediaQueryCursor{
		ViewMode:     query.ViewMode,
		LastID:       query.LastRecord.ID,
		LastCreateAt: query.LastRecord.CreatedAt,
		LastViews:    query.LastRecord.Views,
	}

	searchParameters := types.NewSearchParameters()
	searchParameters.Add(qkeys.LastCursorKey, lastCursor.Encode())

	queryParameters = append(queryParameters, searchParameters.QueryParameter())

	queryString := "?" + strings.Join(queryParameters, "&")

	basePaths := paths.NewHttpPaths()
	basePaths.DownloaderItems += queryString

	extraData := make(map[string]any)
	extraData[items.DisableHTMXEventKey] = true

	pageData := pages.PageFragmentData{
		BasePaths: basePaths,
		Extra:     extraData,
	}

	err := h.templates.Base.ExecuteTemplate(buf, components.ResultShouldLoadHistoryKey, pageData)
	if err != nil {
		return errInternal(err)
	}

	return nil
}
