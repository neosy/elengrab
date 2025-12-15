package grabberh

import (
	"path/filepath"

	htmxvalues "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/values"
	"github.com/valyala/fasthttp"
)

func (h *GrabberHandlers) GetChannelAvatar(ctx *fasthttp.RequestCtx) {
	channelID := ctx.UserValue(channelIDKey).(string)
	if channelID == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("ChannelID is required")
		return
	}

	channelInfo, _ := h.usecases.Downloader.GetYoutubeChannelInfo(ctx, channelID)

	if channelInfo != nil && len(channelInfo.ImageRaw) > 0 {
		ctx.SetContentType(h.mappers.MapImageFormatToContentType(channelInfo.ImageFormat))
		ctx.SetBody(channelInfo.ImageRaw)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")
	defaultAvatarSVG := htmxvalues.IconFileRawByKey(htmxvalues.YoutubeChannelDefaultIconNameKey, iconsDir)

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultAvatarSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
