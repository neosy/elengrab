package handlers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GrabHandler(ctx *fasthttp.RequestCtx) {
	// var pageHasDivItems = cookiePageHasDivItemsKey.compareValue(ctx, "true")

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		err := errorx.New(
			"Authorization error",
			errorx.ArgHttpStatusCode(fasthttp.StatusUnauthorized),
		).Append(err)
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	url := string(ctx.FormValue(formFieldMediaURLKey))
	if url == "" {
		err := errorx.New("URL is required", errorx.ArgHttpStatusCode(fasthttp.StatusBadRequest))
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if err := h.validators.Validate.Var(url, "url"); err != nil {
		err := errorx.New("Invalid URL", errorx.ArgHttpStatusCode(fasthttp.StatusBadRequest))
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	// Read selected quality and format
	formSelectQualityCodec := string(ctx.FormValue(formFieldQualityCodecKey))
	formSelectQualityResolution := string(ctx.FormValue(formFieldQualityResolutionKey))
	formSelectFormat := string(ctx.FormValue(formFieldFormatKey))

	resp, err := h.usecases.Downloader.ScheduleDownload(
		ctx,
		userID,
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
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("the request returned an empty")
		return
	}

	ctx.Response.Header.SetCookie(cookiePageHasDivItemsKey.makeCookie("true", "/", 7*24*60*60))

	ctx.SetStatusCode(fasthttp.StatusOK)
}
