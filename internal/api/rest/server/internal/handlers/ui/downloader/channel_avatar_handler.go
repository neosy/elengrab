package downloader

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetChannelAvatarHandler1(ctx *fasthttp.RequestCtx) {
	channelID, ok := ctx.UserValue(qkeys.ChannelIDKey.String()).(string)
	platform, ok := ctx.UserValue(qkeys.ChannelPlatformKey.String()).(string)
	if !ok || channelID != "" {
		channelInfo, _ := h.downloader.FindChannelInfo(ctx, channelID, platform)

		var image *dtypes.ImageData

		if channelInfo != nil && channelInfo.Image.IsValid() {
			image = (*dtypes.ImageData)(channelInfo.Image)
		}

		if image == nil {
			image, _ = h.downloader.GetSiteImage(ctx, channelInfo.ChannelURL)
		}

		if !image.IsZero() {
			ctx.SetContentType(httpx.ContentTypeByExt(image.Format.String()))
			ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
			ctx.SetBody(image.Raw)
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}
	}

	defaultAvatarSVG := icons.MediaDefaultIcon.FileRaw()

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultAvatarSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
