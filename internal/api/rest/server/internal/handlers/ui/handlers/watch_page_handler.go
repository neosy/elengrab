package handlers

import (
	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) WatchPageByDownloadIDHandler(ctx *fasthttp.RequestCtx) {
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

	h.renderWatchPage(ctx,
		renderWatchPageRequest{
			pageURL:        httppaths.BuildPathMediaItemWatch(downloadID),
			streamPath:     httppaths.BuildPathMediaItemStream(downloadID),
			downloadID:     downloadID,
			showBackButton: true,
		},
	)
}
