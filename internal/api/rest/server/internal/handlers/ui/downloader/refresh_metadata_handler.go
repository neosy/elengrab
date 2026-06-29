package downloader

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) RefreshMetadataByDownloadIDHandler(ctx *fasthttp.RequestCtx) {
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

	downloadID, err := idcodec.DecodeUUIDBase64URL(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	err = h.downloader.ScheduleRefreshMetadata(ctx, authCtx, downloadID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	nfasthttp.WriteResponse(ctx, nil)
}
