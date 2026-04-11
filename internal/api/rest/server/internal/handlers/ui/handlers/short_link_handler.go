package handlers

import (
	"strings"

	"github.com/google/uuid"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ShortLinkHandler(ctx *fasthttp.RequestCtx) {
	shortCode := ctx.UserValue(shortCodeKey).(string)
	if shortCode == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("shortCode is required", fasthttp.StatusBadRequest))
		return
	}

	shortURL, ipAddress, userAgent, referrer := h.extractRequestMeta(ctx)

	link, err := h.linkWeb.ShortLinkClick(ctx, shortURL, ipAddress, userAgent, referrer)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if link == nil {
		nfasthttp.WriteErrorx(
			ctx,
			errorx.NewWithMessage("short link not found.", exceptionx.NOT_FOUND))
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

	if strings.Contains(link.OriginalURL, httppaths.GroupDownloader+httppaths.PathStream+"/") {
		parts := strings.Split(link.OriginalURL, "/")
		if len(parts) == 0 {
			writeError()
			return
		}

		uuidStr := parts[len(parts)-1]
		fileID, err := uuid.Parse(uuidStr)
		if err != nil {
			writeError()
			return
		}

		streamPath := httppaths.BuildPathStreamShortCode(shortCode)
		h.view(ctx, shortURL, streamPath, fileID)
		return
	}

	writeError()
}
