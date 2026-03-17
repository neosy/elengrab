package handlers

import (
	"bytes"
	"fmt"
	"time"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetFilesHistoryHandler(ctx *fasthttp.RequestCtx) {
	var (
		before     = time.Now().UTC()
		bodyBuffer bytes.Buffer
	)

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetBodyString(fmt.Sprintf("Authorization error: %v", err))
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

	// We upload one more line to see if we need to show "Upload more"
	resps, err := h.usecases.Downloader.LoadHistory(ctx, userID, before, loadHistoryLimit+1)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyString("")
		return
	}

	loadNextHistory := len(resps) > loadHistoryLimit

	// If there are more items than the limit, we show only the limited number of items and a "Load more"
	lines := resps
	if len(resps) > loadHistoryLimit {
		lines = resps[:loadHistoryLimit]
	}

	before = time.Time{}

	for i, fileInfo := range lines {
		buf, _, err := h.genRow(fileInfo, false)
		if err != nil || buf == nil {
			continue
		}
		bodyBuffer.Write(buf.Bytes())
		before = fileInfo.CreatedAt

		if loadNextHistory && i == preloadHistoryAfter-1 {
			buf, err := h.genRowShouldLoadHistory(before)
			if err == nil && buf != nil {
				bodyBuffer.Write(buf.Bytes())
			}
		}
	}

	if loadNextHistory {
		buf, err := h.genRowLoadHistory()
		if err == nil && buf != nil {
			bodyBuffer.Write(buf.Bytes())
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}

func (h *DownloaderHandlers) genRowLoadHistory() (*bytes.Buffer, error) {
	dataMap := uivalues.MergeMaps()
	dataMap[uivalues.DisableHTMXEventKey] = true

	var buf bytes.Buffer
	err := h.templates.ExecuteTemplate(&buf, uivalues.GrabResultLoadHistoryHtmlFileName, dataMap)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func (h *DownloaderHandlers) genRowShouldLoadHistory(before time.Time) (*bytes.Buffer, error) {
	if before.IsZero() {
		return nil, nil
	}

	dataMap := uivalues.MergeMaps(uivalues.PathValues)
	dataMap[uivalues.PathItemsHistoryKey] = dataMap[uivalues.PathItemsHistoryKey].(string) + fmt.Sprintf("?before=%s", before.Format(dateFormate))
	dataMap[uivalues.DisableHTMXEventKey] = true
	var buf bytes.Buffer
	err := h.templates.ExecuteTemplate(&buf, uivalues.GrabResultShouldLoadHistoryHtmlFileName, dataMap)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
