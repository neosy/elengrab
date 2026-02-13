package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) RepeatDownloadHandler(ctx *fasthttp.RequestCtx) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetBodyString(fmt.Sprintf("Authorization error: %v", err))
		return
	}

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is required")
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is incorrect")
		return
	}

	fileInfo, err := h.usecases.Downloader.RepeatDownload(ctx, userID, fileId)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	buf, httpStatus, err := h.genRow(fileInfo, false)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}
	if httpStatus == fasthttp.StatusNoContent {
		ctx.SetStatusCode(httpStatus)
		ctx.Response.Header.Set("HX-Trigger", "no-op")
		return
	}

	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
