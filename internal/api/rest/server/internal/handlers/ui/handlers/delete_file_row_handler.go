package handlers

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/dto"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) DeleteRowHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := policy.ResolveUserOrFallback(ctx, h.appMode)
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

	dataResp := dto.GrabResponse{
		GuestCreated: ctxUser.GuestCreated,
	}

	dataJSON, _ := json.Marshal(dataResp)
	if dataJSON != nil {
		ctx.Response.Header.Set("HX-Trigger", string(dataJSON))
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
