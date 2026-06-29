package static

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) ThumbnailHandler(ctx *fasthttp.RequestCtx) {
	thumbnailIdStr, ok := ctx.UserValue(thumbnailIdKey).(string)
	if !ok || thumbnailIdStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrThumbnailIdIsRequired)
		return
	}

	thumbnailID, err := idcodec.DecodeUUIDBase64URL(thumbnailIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrThumbnailIdIsIncorrect.Wrap(err))
		return
	}

	thumbnail, err := h.thumbnail.LoadByThumbID(ctx, thumbnailID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	ctx.SetContentType(httpx.ContentTypeByExt(thumbnail.Format.String()))
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody(thumbnail.ImageRaw)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
