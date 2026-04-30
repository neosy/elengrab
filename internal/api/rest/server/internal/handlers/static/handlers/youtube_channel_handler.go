package handlers

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) YoutubeChannelHandler(ctx *fasthttp.RequestCtx) {
	channelID := ctx.UserValue(ChannelIdKey).(string)
	if channelID == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrChannelIsRequired)
		return
	}

	channel, err := h.downloader.GetYoutubeChannelInfo(ctx, channelID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	ctx.SetContentType(httpx.ContentTypeByExt(channel.ImageFormat.String()))
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody(channel.ImageRaw)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
