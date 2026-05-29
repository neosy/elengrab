package handlers

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/policy"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaHistoryHandler(ctx *fasthttp.RequestCtx) {
	var before = time.Now().UTC()

	ctxUser := policy.ResolveUserOrAnonym(ctx)

	beforeStr := string(ctx.QueryArgs().Peek(beforeKey))
	if beforeStr != "" {
		var err error
		before, err = time.Parse(dateFormate, beforeStr)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBodyString("")
			return
		}
	}

	filters := parseFilters(ctx)

	var bodyBuffer bytes.Buffer
	err := h.getDownloadsHistory(ctx, &bodyBuffer, *ctxUser, before, filters)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyString("")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}

func (h *DownloaderHandlers) getDownloadsHistory(
	ctx context.Context,
	buf *bytes.Buffer,
	userCtx dauth.UserContext,
	before time.Time,
	filters requestFilters,
) error {
	var filterByTitle string
	if filters != nil {
		filterByTitle = filters["title"]
	}

	// We upload one more line to see if we need to show "Upload more"
	resps, err := h.downloader.LoadHistory(ctx, userCtx, before, loadHistoryLimit+1, filterByTitle)
	if err != nil {
		return err
	}

	if len(resps) == 0 {
		return nil
	}

	loadNextHistory := len(resps) > loadHistoryLimit

	// If there are more items than the limit, we show only the limited number of items and a "Load more"
	lines := resps
	if len(resps) > loadHistoryLimit {
		lines = resps[:loadHistoryLimit]
	}

	before = lines[len(lines)-1].CreatedAt

	for i, downloadInfo := range lines {
		row := h.renderMediaItemRow(downloadInfo, false)
		if row.err != nil {
			h.logger.Warn(
				"Failed to generate row",
				"error", err,
			)
			continue
		}

		// Load template
		tmpl, err := h.templates.Clone()
		if err != nil {
			h.logger.Warn(
				"Failed to clone templates",
				"error", err,
			)
			continue
		}

		err = tmpl.ExecuteTemplate(buf, components.ResultRowStatusKey, row.data)
		if err != nil {
			h.logger.Warn(
				"Failed to execute template",
				"name", components.ResultRowStatusKey,
				"error", err,
			)
			continue
		}

		if loadNextHistory && i == preloadHistoryAfter-1 {
			h.renderRowShouldLoadHistory(buf, before, filters)
		}
	}

	if loadNextHistory {
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

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		return errInternal(err)
	}

	err = tmpl.ExecuteTemplate(buf, components.ResultLoadHistory, pageData)
	if err != nil {
		return errInternal(err)
	}

	return nil
}

func (h *DownloaderHandlers) renderRowShouldLoadHistory(
	buf *bytes.Buffer,
	before time.Time,
	filters requestFilters,
) error {
	if before.IsZero() {
		return nil
	}

	queryString := fmt.Sprintf("?before=%s", before.Format(dateFormate))

	if filters != nil {
		filterByTitle := filters[filterByTitleKey]
		queryString += fmt.Sprintf("&filter[%s]=%s", filterByTitleKey, filterByTitle)
	}

	basePaths := paths.NewPaths()
	basePaths.DownloaderHistory += queryString

	extraData := make(map[string]any)
	extraData[items.DisableHTMXEventKey] = true

	pageData := pages.PageFragmentData{
		BasePaths: basePaths,
		Extra:     extraData,
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		return errInternal(err)
	}

	err = tmpl.ExecuteTemplate(buf, components.ResultShouldLoadHistoryKey, pageData)
	if err != nil {
		return errInternal(err)
	}

	return nil
}
