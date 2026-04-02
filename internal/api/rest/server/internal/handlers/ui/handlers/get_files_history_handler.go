package handlers

import (
	"bytes"
	"context"
	"fmt"
	"time"

	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetFilesHistoryHandler(ctx *fasthttp.RequestCtx) {
	var before = time.Now().UTC()

	ctxUser := authmw.UserFromContext(ctx)
	if ctxUser == nil {
		// anonymous
		ctxUser = dauth.UserContextAnonymous()
	}

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
	err := h.getFilesHistory(ctx, &bodyBuffer, *ctxUser, before, filters)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyString("")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}

func (h *DownloaderHandlers) getFilesHistory(
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

	for i, fileInfo := range lines {
		row := h.genRow(fileInfo, false)
		if row.err != nil {
			continue
		}

		// Load template
		tmpl, err := h.templates.Clone()
		if err != nil {
			continue
		}

		err = tmpl.ExecuteTemplate(buf, row.templateName, row.data)
		if err != nil {
			continue
		}

		if loadNextHistory && i == preloadHistoryAfter-1 {
			h.genRowShouldLoadHistory(buf, before, filters)
		}
	}

	if loadNextHistory {
		h.genRowLoadHistory(buf)
	}

	return nil
}

func (h *DownloaderHandlers) genRowLoadHistory(buf *bytes.Buffer) error {
	dataMap := uivalues.MergeMaps()
	dataMap[uivalues.DisableHTMXEventKey] = true

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		return errInternal(err)
	}

	err = tmpl.ExecuteTemplate(buf, uivalues.ComponentResultLoadHistory, dataMap)
	if err != nil {
		return errInternal(err)
	}

	return nil
}

func (h *DownloaderHandlers) genRowShouldLoadHistory(
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

	dataMap := uivalues.MergeMaps(uivalues.PathValues)
	dataMap[uivalues.PathItemsHistoryKey] = dataMap[uivalues.PathItemsHistoryKey].(string) + queryString
	dataMap[uivalues.DisableHTMXEventKey] = true

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		return errInternal(err)
	}

	err = tmpl.ExecuteTemplate(buf, uivalues.ComponentResultShouldLoadHistoryKey, dataMap)
	if err != nil {
		return errInternal(err)
	}

	return nil
}
