package downloader

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) EditMediaPageByDownloadIDHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		fasthttpx.WriteHTMLHeadResponse(ctx)
		return
	}

	authCtx, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

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

	h.renderEditMediaPage(ctx,
		renderEditMediaPageRequest{
			pageURL:    httppaths.BuildMediaItemPath(downloadID),
			downloadID: downloadID,
			authCtx:    authCtx,
		},
	)
}
