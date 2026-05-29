package handlers

import (
	"strings"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader_handlers/dto"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/policy"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ImportMediaByURLHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	url := string(ctx.FormValue(formFieldMediaURLKey))
	url = strings.TrimSpace(url)
	if url == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrURLIsRequired)
		return
	}

	if err := h.validators.Validate.Var(url, "url"); err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrInvalidURL.Wrap(err))
		return
	}

	// Read selected quality and format
	formSelectQualityCodec := string(ctx.FormValue(formFieldQualityCodecKey))
	formSelectQualityResolution := string(ctx.FormValue(formFieldQualityResolutionKey))
	formSelectFormat := string(ctx.FormValue(formFieldFormatKey))

	resp, err := h.downloader.ScheduleDownload(
		ctx,
		*ctxUser,
		url,
		&ddownload.DownloadOptions{
			FormatType:      h.mappers.MapFormatType(formSelectQualityCodec, formSelectFormat),
			VideoFormat:     h.mappers.MapVideoFormat(formSelectQualityCodec, formSelectFormat),
			VideoCodec:      h.mappers.MapVideoCodec(formSelectQualityCodec),
			VideoResolution: h.mappers.MapVideoResolution(formSelectQualityResolution),
			AudioFormat:     h.mappers.MapAudioFormat(formSelectFormat),
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

	dataResp := dto.GrabResponse{
		GuestCreated: ctxUser.GuestCreated,
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	nfasthttp.WriteResponse(ctx, dataResp)
}
