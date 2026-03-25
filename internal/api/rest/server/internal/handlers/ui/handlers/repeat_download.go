package handlers

import (
	"bytes"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) RepeatDownloadHandler(ctx *fasthttp.RequestCtx) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("authorization error", fasthttp.StatusUnauthorized, err))
		return
	}

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("fileId is required", fasthttp.StatusBadRequest))
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("fileId is incorrect", fasthttp.StatusBadRequest))
		return
	}

	fileInfo, err := h.usecases.Downloader.RepeatDownload(ctx, userID, fileId)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	row := h.genRow(fileInfo, false)
	if row.err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}
	if row.httpStatus == fasthttp.StatusNoContent {
		ctx.SetStatusCode(row.httpStatus)
		ctx.Response.Header.Set("HX-Trigger", "no-op")
		return
	}

	var buf bytes.Buffer
	err = h.templates.ExecuteTemplate(&buf, row.templateName, row.data)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err))
		return
	}

	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
