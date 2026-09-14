package downloader

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetChannelAvatarHandler(ctx *fasthttp.RequestCtx) {
	channelID, ok := ctx.UserValue(qkeys.ChannelIDKey.String()).(string)
	platform, ok := ctx.UserValue(qkeys.ChannelPlatformKey.String()).(string)
	if !ok || channelID != "" {
		channelInfo, _ := h.downloader.FindChannelInfo(ctx, channelID, platform)

		if channelInfo != nil && channelInfo.Image.IsValid() {
			ctx.SetContentType(httpx.ContentTypeByExt(channelInfo.Image.Format.String()))
			ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
			ctx.SetBody(channelInfo.Image.Raw)
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
