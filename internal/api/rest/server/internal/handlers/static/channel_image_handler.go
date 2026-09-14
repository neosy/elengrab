package static

import (
	"net/http"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/pkg/errorx"
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

	if !channel.HasImage() {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("channel image not found", http.StatusNotFound))
		return
	}

	ctx.SetContentType(httpx.ContentTypeByExt(channel.Image.Format.String()))
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody(channel.Image.Raw)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
