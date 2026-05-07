package handlers

import (
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *V1Handlers) GetChannelByID(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()

	channelID := string(args.Peek(channelIDKey))
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
