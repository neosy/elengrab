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

	etag := thumbnail.ETag()

	ctx.Response.Header.Set(fasthttp.HeaderETag, etag)
	ctx.Response.Header.Set(fasthttp.HeaderCacheControl, "public, max-age=86400")
	ctx.Response.Header.SetLastModified(thumbnail.UpdatedAt)

	// Conditional request: the client already has this version
	if string(ctx.Request.Header.Peek(fasthttp.HeaderIfNoneMatch)) == etag {
		ctx.SetStatusCode(fasthttp.StatusNotModified)
		return
	}

	ctx.SetContentType(httpx.ContentTypeByExt(thumbnail.Format.String()))
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentLength(len(thumbnail.ImageRaw))

	if ctx.IsHead() {
		return
	}

	ctx.SetBody(thumbnail.ImageRaw)
}
