package handlers

import (
	"strings"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/dto"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) CreateMediaShareLinkHandler(ctx *fasthttp.RequestCtx) {
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

	url, err := h.linkWeb.CreateShortLink(
		ctx,
		strings.TrimSuffix(h.baseURL, "/")+httppaths.BuildPathMediaItemWatch(downloadID),
	)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	resp := dto.GetDownloadShareLinkResponse{
		DownloadID: downloadID.String(),
		URL:        url,
	}

	nfasthttp.WriteResponse(ctx, resp)
}
