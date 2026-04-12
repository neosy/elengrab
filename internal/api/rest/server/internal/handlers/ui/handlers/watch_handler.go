package handlers

import (
	"github.com/google/uuid"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) FileWatchHandler(ctx *fasthttp.RequestCtx) {
	// Get user ID from context
	ctxUser := authmw.UserFromContext(ctx)
	if ctxUser == nil {
		// anonymous
		ctxUser = dauth.UserContextAnonymous()
	}

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("fileId is required", fasthttp.StatusBadRequest))
		return
	}

	fileID, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("fileId is incorrect", fasthttp.StatusBadRequest, err))
		return
	}

	h.watch(
		ctx,
		httppaths.BuildPathFileWatch(fileID),
		httppaths.BuildPathFileStream(fileID),
		fileID,
		true,
	)
}
