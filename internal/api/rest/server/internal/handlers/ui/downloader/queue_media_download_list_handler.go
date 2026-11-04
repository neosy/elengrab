package downloader

import (
	"encoding/json"
	"net/http"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) QueueMediaDownloadListHandler(ctx *fasthttp.RequestCtx) {
	authCtx, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var receivedReq dto.QueueMediaDownloadListRequest
	if err := json.Unmarshal(ctx.PostBody(), &receivedReq); err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("Invalid JSON request body", http.StatusBadRequest, err))
		return
	}

	if err := h.validators.Validate.Struct(receivedReq); err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	if err := receivedReq.Validate(); err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	req, err := h.mappers.MapQueueMediaDownloadListRequestToScheduleDownloadListRequest(receivedReq)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	scheduleResp, err := h.downloader.ScheduleDownloadList(ctx, authCtx, req)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if scheduleResp == nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrEmptyResponse)
		return
	}

	dataResponse, err := h.mappers.MapScheduleDownloadListResponseToQueueMediaDownloadListResponse(scheduleResp)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	nfasthttp.WriteResponse(ctx, dataResponse)
}
