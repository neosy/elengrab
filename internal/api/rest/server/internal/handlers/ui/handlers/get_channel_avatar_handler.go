package handlers

import (
	"path/filepath"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetChannelAvatarHandler(ctx *fasthttp.RequestCtx) {
	channelID := ctx.UserValue(channelIDKey).(string)
	if channelID == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("channelID is required", fasthttp.StatusBadRequest))
		return
	}

	if channelID != channelIDValueNone {
		channelInfo, _ := h.usecases.Downloader.FindYoutubeChannelInfo(ctx, channelID)

		if channelInfo != nil && len(channelInfo.ImageRaw) > 0 {
			ctx.SetContentType(h.mappers.MapImageExtToContentType(channelInfo.ImageFormat))
			ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
			ctx.SetBody(channelInfo.ImageRaw)
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}
	}

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")
	defaultAvatarSVG := uivalues.IconFileRawByKey(uivalues.MediaDefaultIconNameKey, iconsDir)

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultAvatarSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
