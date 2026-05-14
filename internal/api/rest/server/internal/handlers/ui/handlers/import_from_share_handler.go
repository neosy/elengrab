package handlers

import (
	"strings"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ImportFromShareHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	args := ctx.QueryArgs()

	url := string(args.Peek(urlKey))
	if url == "" {
		url = string(args.Peek(textKey))
	}

	url = strings.TrimSpace(url)
	if url == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrURLIsRequired)
		return
	}

	if err := h.validators.Validate.Var(url, "url"); err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrInvalidURL.Wrap(err))
		return
	}

	resp, err := h.downloader.ScheduleDownload(
		ctx,
		*ctxUser,
		url,
		&ddownload.DownloadOptions{
			FormatType:      dtypes.FormatTypeVideoAudio,
			VideoCodec:      new(dtypes.VideoCodecH264),
			VideoResolution: new(dtypes.VideoResolution1080p),
			VideoFormat:     new(dtypes.VideoFormatAuto),
		},
	)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if resp == nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrEmptyResponse)
		return
	}

	ctx.Redirect(httppaths.PathIndex, fasthttp.StatusFound)
}
