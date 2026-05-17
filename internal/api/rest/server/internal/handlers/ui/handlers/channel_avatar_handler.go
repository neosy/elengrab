package handlers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetChannelAvatarHandler(ctx *fasthttp.RequestCtx) {
	channelID, ok := ctx.UserValue(channelIDKey).(string)
	if !ok || channelID != "" {
		channelInfo, _ := h.downloader.FindYoutubeChannelInfo(ctx, channelID)

		if channelInfo != nil && len(channelInfo.ImageRaw) > 0 {
			ctx.SetContentType(httpx.ContentTypeByExt(channelInfo.ImageFormat.String()))
			ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
			ctx.SetBody(channelInfo.ImageRaw)
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}
	}

	defaultAvatarSVG := icons.FileRawByKey(icons.MediaDefaultIconNameKey, h.assetFolders.Icons())

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultAvatarSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
