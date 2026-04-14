package handlers

import (
	"bytes"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetFileRowHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)

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

	fileInfo, err := h.downloader.GetFileInfo(ctx, *ctxUser, fileId)
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

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, row.templateName, row.data)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
