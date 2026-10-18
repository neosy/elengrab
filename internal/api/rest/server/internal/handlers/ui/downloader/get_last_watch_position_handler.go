package downloader

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetLastWatchPositionHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		fasthttpx.WriteHTMLHeadResponse(ctx)
		return
	}

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

	lastWatchPosition, err := h.downloader.GetLastWatchPosition(ctx, authCtx, downloadID)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	resp := dto.GetLastWatchPositionResponse{
		DownloadID: idcodec.EncodeUUIDBase64URL(downloadID),
		Position:   uint32(lastWatchPosition.Milliseconds()),
	}

	fasthttpx.WriteResponse(ctx, resp)
}
