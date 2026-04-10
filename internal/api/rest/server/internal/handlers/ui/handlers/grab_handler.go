package handlers

import (
	"strings"

	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GrabHandler(ctx *fasthttp.RequestCtx) {
	// var pageHasDivItems = cookiePageHasDivItemsKey.compareValue(ctx, "true")

	ctxUser, err := authmw.EnsureUserFromContext(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	url := string(ctx.FormValue(formFieldMediaURLKey))
	if url == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("URL is required", fasthttp.StatusBadRequest))
		return
	}

	url = strings.TrimSpace(url)

	if err := h.validators.Validate.Var(url, "url"); err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("invalid URL", fasthttp.StatusBadRequest, err))
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
		nfasthttp.WriteErrorx(
			ctx,
			errorx.NewHTTP(
				"the request returned an empty response",
				fasthttp.StatusInternalServerError,
			),
		)
		return
	}

	ctx.Response.Header.SetCookie(cookiePageHasDivItemsKey.makeCookie("true", "/", 7*24*60*60))

	ctx.SetStatusCode(fasthttp.StatusCreated)
}
