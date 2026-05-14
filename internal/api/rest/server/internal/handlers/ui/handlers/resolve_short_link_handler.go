package handlers

import (
	"github.com/google/uuid"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ResolveShortLinkHandler(ctx *fasthttp.RequestCtx) {
	shortCode, ok := ctx.UserValue(shortCodeKey).(string)
	if !ok || shortCode == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("shortCode is required", fasthttp.StatusBadRequest))
		return
	}

	shortURL, ipAddress, userAgent, referrer := h.extractRequestMeta(ctx)

	link, err := h.linkWeb.ShortLinkClick(ctx, shortURL, ipAddress, userAgent, referrer)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if link == nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewWithMessage("short link not found.", exceptionx.NOT_FOUND))
		return
	}

	writeError := func() {
		nfasthttp.WriteErrorx(
			ctx,
			errorx.Errorf(
				"failed short link %v, originalURL: %v", link.ShortURL, link.OriginalURL,
				exceptionx.WRONG_DATA,
				errorx.WithErrorMessage("Short link is invalid."),
			))
	}

	if downloadID := stripUUIDFromPath(link.OriginalURL); downloadID != uuid.Nil {
		streamPath := httppaths.BuildPathStreamShortCode(shortCode)
		h.renderWatchPage(ctx,
			renderWatchPageRequest{
				pageURL:        shortURL,
				streamPath:     streamPath,
				downloadID:     downloadID,
				showBackButton: false,
			},
		)
		return
	}

	writeError()
}
