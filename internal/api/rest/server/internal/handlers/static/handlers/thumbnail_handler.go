package handlers

import (
	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) ThumbnailHandler(ctx *fasthttp.RequestCtx) {
	thumbnailIdStr := ctx.UserValue(thumbnailIdKey).(string)
	if thumbnailIdStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrThumbnailIdIsRequired)
		return
	}

	thumbnailID, err := uuid.Parse(thumbnailIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrThumbnailIdIsIncorrect.Wrap(err))
		return
	}

	thumbnail, err := h.thumbnail.GetByThumbID(ctx, thumbnailID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	ctx.SetContentType(httpx.ContentTypeByExt(thumbnail.Format.String()))
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody(thumbnail.ImageRaw)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
