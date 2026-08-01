package downloader

import (
	"encoding/json"
	"net/http"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemWatchTrackingHandler(ctx *fasthttp.RequestCtx) {
	// Get user ID from context
	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadIDStr, ok := ctx.UserValue(downloadIDKey).(string)
	if !ok || downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := idcodec.DecodeUUIDBase64URL(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	var receivedReq dto.MediaItemWatchTrackingRequest
	if err := json.Unmarshal(ctx.PostBody(), &receivedReq); err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("Invalid JSON request body", http.StatusBadRequest, err))
		return
	}

	if err := h.validators.Validate.Struct(receivedReq); err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	req, err := h.mappers.MapMediaItemWatchTrackingRequestToUsecase(authCtx, downloadID, receivedReq)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	err = h.downloader.TrackMediaWatchEvent(ctx, authCtx, req)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	nfasthttp.WriteResponse(ctx, nil)
}
