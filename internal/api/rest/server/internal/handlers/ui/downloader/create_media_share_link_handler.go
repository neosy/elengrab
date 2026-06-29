package downloader

import (
	"strings"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) CreateMediaShareLinkHandler(ctx *fasthttp.RequestCtx) {
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
