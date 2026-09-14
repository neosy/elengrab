package apiv1

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *V1Handlers) GetChannelByID(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()

	channelID := string(args.Peek(channelIDKey))
	if channelID == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrChannelIsRequired)
		return
	}

	platform := string(args.Peek(platformKey))
	if platform == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrChannelPlatformIsRequired)
		return
	}

	channel, err := h.downloader.GetChannelInfo(ctx, channelID, platform)
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
