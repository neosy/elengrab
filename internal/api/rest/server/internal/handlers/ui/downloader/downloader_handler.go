package downloaderh

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/valyala/fasthttp"
)

type grabResultData struct {
	PathFileRow      string
	YoutubeChannelID string
	YoutubeTitle     string
	YoutubeURL       string
	FileSize         string
	Format           string
	DownloadURL      string
}

func (h *DownloaderHandlers) GrabHandler(ctx *fasthttp.RequestCtx) {
	var pageHasDivItems = cookiePageHasDivItemsKey.compareValue(ctx, "true")

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetBodyString(fmt.Sprintf("Authorization error: %v", err))
		return
	}

	url := string(ctx.FormValue(formFieldYouTubeURLKey))
	if url == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("URL is required")
		return
	}

	if err := h.validators.Validate.Var(url, "url"); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Invalid URL")
		return
	}

	// Read selected quality and format
	formSelectQualityCodec := string(ctx.FormValue(formFieldQualityCodecKey))
	formSelectQualityResolution := string(ctx.FormValue(formFieldQualityResolutionKey))
	formSelectFormat := string(ctx.FormValue(formFieldFormatKey))

	var (
		data = grabResultData{
			YoutubeURL:       url,
			YoutubeChannelID: channelIDValueNone,
		}
	)

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
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(fmt.Sprintf("Internal error: %v", err))
		return
	}

	if resp == nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("the request returned an empty")
		return
	}

	ctx.Response.Header.SetCookie(cookiePageHasDivItemsKey.makeCookie("true", "/", 7*24*60*60))

	// Updating the cache
	{
		cacheRow.mu.Lock()
		cacheRow.data[resp.FileId] = cacheRowEntry{
			youtubeTitle: resp.YoutubeTitle,
			Format:       resp.Format,
			Updated:      time.Now(),
		}
		cacheRow.mu.Unlock()
	}

	data.YoutubeTitle = url
	data.PathFileRow = httppaths.BuildPathFileRow(resp.FileId)
	data.FileSize = "-"
	data.Format = "-"
	// Set URL for download endpoint
	data.DownloadURL = fmt.Sprintf("%s?file=%s", httppaths.GroupDownloader+httppaths.PathDownload, resp.FileId)

	dataMap := uivalues.MergeMaps(
		uivalues.PathValues,
		uivalues.IconFileNames(),
		uivalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	dataMap[uivalues.GrabResultStatusIconNameKey] = uivalues.FrabResultStatusIconFileName(resp.Status)
	dataMap[uivalues.GrabResultItemStatusHtmlKey] = template.HTML(
		uivalues.GrabResultStatusIconSvgRaw(resp.Status, iconsDir))
	dataMap[uivalues.GrabResultItemDeleteIconKey] = template.HTML(
		uivalues.IconFileRaw(uivalues.IconFileName(uivalues.DownloadDeleteIconNameKey), iconsDir))
	dataMap[uivalues.IsItemHTMXOptionRepeatKey] = true
	dataMap[uivalues.PageHasDivItemsKey] = pageHasDivItems
	dataMap[uivalues.ItemFadeKey] = "fade-in"
	dataMap[uivalues.GrabResultItemStatusTextKey] = ""
	dataMap[uivalues.ResultYoutubeUrlFadeKey] = ""
	dataMap[uivalues.ResultSizeFadeKey] = ""
	dataMap[uivalues.ResultFormatFadeKey] = ""

	var buf bytes.Buffer
	h.templates.ExecuteTemplate(&buf, uivalues.GrabResultNewItemHtmlFileName, dataMap)
	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
