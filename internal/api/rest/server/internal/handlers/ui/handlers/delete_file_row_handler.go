package handlers

import (
	"github.com/google/uuid"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) DeleteFileRowHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := authmw.EnsureUserFromContext(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("fileId is required", fasthttp.StatusBadRequest))
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("fileId is incorrect", fasthttp.StatusBadRequest, err))
		return
	}

	err = h.downloader.DeleteDownload(ctx, *ctxUser, fileId)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
