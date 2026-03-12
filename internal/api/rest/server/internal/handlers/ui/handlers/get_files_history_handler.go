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

	resps, err := h.usecases.Downloader.LoadHistory(ctx, userID, before, loadHistoryLimit)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyString("")
		return
	}

	before = time.Time{}
	for _, fileInfo := range resps {
		buf, _, err := h.genRow(fileInfo, false)
		if err != nil || buf == nil {
			continue
		}
		bodyBuffer.Write(buf.Bytes())
		before = fileInfo.CreatedAt
	}

	buf, _, err := h.genRowLoadHistory(before)
	if err == nil && buf != nil {
		bodyBuffer.Write(buf.Bytes())
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}

func (h *DownloaderHandlers) genRowLoadHistory(before time.Time) (*bytes.Buffer, int, error) {
	if before.IsZero() {
		return nil, fasthttp.StatusOK, nil
	}

	dataMap := uivalues.MergeMaps(uivalues.PathValues)
	dataMap[uivalues.PathItemsHistoryKey] = dataMap[uivalues.PathItemsHistoryKey].(string) + fmt.Sprintf("?before=%s", before.Format(dateFormate))

	var buf bytes.Buffer
	err := h.templates.ExecuteTemplate(&buf, uivalues.GrabResultLoadHistoryHtmlFileName, dataMap)
	if err != nil {
		return nil, fasthttp.StatusInternalServerError, err
	}

	return &buf, fasthttp.StatusOK, nil
}
