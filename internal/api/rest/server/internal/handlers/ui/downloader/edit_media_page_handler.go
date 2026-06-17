package downloader

import (
	"mime"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) EditMediaPageByDownloadIDHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		ctx.SetContentType(mime.TypeByExtension(".html"))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

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

	h.renderEditMediaPage(ctx,
		renderEditMediaPageRequest{
			pageURL:    httppaths.BuildPathMediaItemWatch(downloadID),
			downloadID: downloadID,
			authCtx:    authCtx,
		},
	)
}
