package downloader

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) PatchMediaByDownloadIDHandler(ctx *fasthttp.RequestCtx) {
	authCtx, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	downloadIDStr, ok := ctx.UserValue(downloadIDKey).(string)
	if !ok || downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := uuid.Parse(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	var receivedReq dto.PatchMediaByDownloadIDRequest
	if err := json.Unmarshal(ctx.PostBody(), &receivedReq); err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("Invalid JSON request body", http.StatusBadRequest, err))
		return
	}

	req, err := h.mappers.MapPatchMediaRequestToUsecase(downloadID, receivedReq)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	err = h.downloader.PatchMediaDownload(ctx, authCtx, req)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	nfasthttp.WriteResponse(ctx, nil)
}
