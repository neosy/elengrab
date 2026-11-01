package static

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) ChannelImageHandler(ctx *fasthttp.RequestCtx) {
	externalID, ok := ctx.UserValue(ExternalIdKey).(string)
	if !ok || externalID == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrChannelIsRequired)
		return
	}

	platform, ok := ctx.UserValue(PlatformKey).(string)
	if !ok || externalID == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrChannelPlatformIsRequired)
		return
	}

	channel, err := h.downloader.GetChannelInfo(ctx, externalID, platform)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var image *dtypes.ImageData

	if channel.HasImage() {
		image = (*dtypes.ImageData)(channel.Image)
	}

	if image.IsZero() && channel.SiteURL() != "" {
		image, _ = h.downloader.GetSiteImage(ctx, channel.SiteURL())
	}

	if image.IsZero() {
		defaultAvatarSVG := icons.MediaDefaultIcon.FileRaw()

		ctx.SetContentType("image/svg+xml")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
		ctx.Response.Header.SetContentLength(len(defaultAvatarSVG))

		if ctx.IsHead() {
			return
		}

		ctx.SetBody([]byte(defaultAvatarSVG))

		return
	}

	ctx.SetContentType(httpx.ContentTypeByExt(image.Format.String()))
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.Response.Header.SetContentLength(len(image.Raw))

	if ctx.IsHead() {
		return
	}

	ctx.SetBody(image.Raw)
}
