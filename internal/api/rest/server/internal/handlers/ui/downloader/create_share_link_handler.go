package downloader

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) CreateShareLinkHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadIDStr, ok := ctx.UserValue(qkeys.DownloadIDKey.String()).(string)
	if !ok || downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := idcodec.DecodeUUIDBase64URL(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	err = h.downloader.CheckDownloadVisibilityAccess(ctx, authCtx, downloadID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	url, err := h.linkWeb.CreateShortLink(
		ctx,
		h.buildMediaWatchURL(downloadID),
	)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	downloadChanged := ucdto.MediaDownloadChanged{
		DownloadID: downloadID,
	}
	downloadChanged.MarkShareLinkChanges()
	h.downloader.NotifyDownloadChanged(ctx, downloadChanged)

	resp := dto.GetShareLinkResponse{
		DownloadID: idcodec.EncodeUUIDBase64URL(downloadID),
		URL:        url,
	}

	nfasthttp.WriteResponse(ctx, resp)
}
