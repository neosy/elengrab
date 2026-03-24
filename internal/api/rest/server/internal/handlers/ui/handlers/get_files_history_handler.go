package handlers

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetFilesHistoryHandler(ctx *fasthttp.RequestCtx) {
	var before = time.Now().UTC()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("authorization error", fasthttp.StatusUnauthorized, err))
		return
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
	err = h.getFilesHistory(ctx, &bodyBuffer, userID, before, filters)
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
	userID uuid.UUID,
	before time.Time,
	filters requestFilters,
) error {
	var filterByTitle string
	if filters != nil {
		filterByTitle = filters["title"]
	}

	// We upload one more line to see if we need to show "Upload more"
	resps, err := h.usecases.Downloader.LoadHistory(ctx, userID, before, loadHistoryLimit+1, filterByTitle)
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

		err = h.templates.ExecuteTemplate(buf, row.templateName, row.data)
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

	err := h.templates.ExecuteTemplate(buf, uivalues.GrabResultLoadHistoryHtmlFileName, dataMap)
	if err != nil {
		return err
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
	err := h.templates.ExecuteTemplate(buf, uivalues.GrabResultShouldLoadHistoryHtmlFileName, dataMap)
	if err != nil {
		return err
	}

	return nil
}
