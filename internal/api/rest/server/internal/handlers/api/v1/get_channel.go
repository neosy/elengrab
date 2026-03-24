package handlers

import (
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *V1Handlers) GetChannelByID(ctx *fasthttp.RequestCtx) {
	channelID := ctx.UserValue(channelIDKey).(string)
	if channelID == "" {
		nfasthttp.WriteErrorx(ctx, errorx.New("channelID is required", exceptionx.WRONG_DATA))
		return
	}

	channel, err := h.usecases.Downloader.GetYoutubeChannelInfo(ctx, channelID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	resp, err := h.mappers.MapChannelDomainToResponse(channel)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	nfasthttp.WriteResponse(ctx, resp)
}
