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
	channelID, ok := ctx.UserValue(ChannelIdKey).(string)
	if !ok || channelID == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrChannelIsRequired)
		return
	}

	platform, ok := ctx.UserValue(PlatformKey).(string)
	if !ok || channelID == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrChannelPlatformIsRequired)
		return
	}

	channel, err := h.downloader.GetChannelInfo(ctx, channelID, platform)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var image *dtypes.ImageData

	if channel.HasImage() {
		image = (*dtypes.ImageData)(channel.Image)
	}

	if image.IsZero() {
		image, _ = h.downloader.GetSiteImage(ctx, channel.ChannelURL)
	}

	if !image.IsZero() {
		ctx.SetContentType(httpx.ContentTypeByExt(image.Format.String()))
		ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
		ctx.SetBody(image.Raw)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	defaultAvatarSVG := icons.MediaDefaultIcon.FileRaw()

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultAvatarSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
