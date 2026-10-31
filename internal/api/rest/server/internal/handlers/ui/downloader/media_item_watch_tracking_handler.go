package downloader

import (
	"encoding/json"
	"net/http"
	"uuid"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemWatchTrackingHandler(ctx *fasthttp.RequestCtx) {
	// Get user ID from context
	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadIDStr, ok := ctx.UserValue(qkeys.DownloadIDKey.String()).(string)
	if !ok || downloadIDStr == "" {
		fasthttpx.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := idcodec.DecodeUUIDBase64URL(downloadIDStr)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	h.watchTracking(ctx, authCtx, downloadID, false)
}

func (h *DownloaderHandlers) ShortLinkWatchTrackingHandler(ctx *fasthttp.RequestCtx) {
	// Get user ID from context
	authCtx := policy.ResolveUserOrAnonym(ctx)

	_, downloadID, err := h.resolveShortLinkToDownloadID(ctx, false)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	h.watchTracking(ctx, authCtx, downloadID, true)
}

func (h *DownloaderHandlers) watchTracking(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
	unrestricted bool,

) {
	var receivedReq dto.MediaItemWatchTrackingRequest
	if err := json.Unmarshal(ctx.PostBody(), &receivedReq); err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewHTTPMessage("Invalid JSON request body", http.StatusBadRequest, err))
		return
	}

	if err := h.validators.Validate.Struct(receivedReq); err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	req, err := h.mappers.MapMediaItemWatchTrackingRequestToUsecase(authCtx, downloadID, receivedReq)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	if !unrestricted {
		err = h.downloader.TrackMediaWatchEvent(ctx, authCtx, req)
		if err != nil {
			fasthttpx.WriteErrorx(ctx, err)
			return
		}
	} else {
		err = h.downloader.TrackMediaWatchEventUnrestricted(ctx, authCtx, req)
		if err != nil {
			fasthttpx.WriteErrorx(ctx, err)
			return
		}
	}

	fasthttpx.WriteResponse(ctx, nil)
}
