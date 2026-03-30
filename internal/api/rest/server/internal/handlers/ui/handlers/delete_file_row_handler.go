package handlers

import (
	"github.com/google/uuid"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) DeleteFileRowHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := authmw.EnsureUserFromContext(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(
			ctx,
			errorx.Errorf("authorization error: %w", err, errorx.HttpStatusArg(fasthttp.StatusUnauthorized)),
		)
		return
	}

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("fileId is required", fasthttp.StatusBadRequest))
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("fileId is incorrect", fasthttp.StatusBadRequest, err))
		return
	}

	err = h.downloader.DeleteDownload(ctx, *ctxUser, fileId)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
