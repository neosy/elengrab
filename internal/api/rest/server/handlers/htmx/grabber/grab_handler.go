package grabberh

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	htmxvalues "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/valyala/fasthttp"
)

type grabResultData struct {
	PathFileRow  string
	YoutubeTitle string
	YoutubeURL   string
	FileSize     string
	Format       string
	DownloadURL  string
}

var (
	formatTypeMap = map[string]dtypes.FormatType{
		"video_orig":     dtypes.FormatTypeVideoAudio,
		"video_mp4_orig": dtypes.FormatTypeVideoAudio,
		"video_mp4_h264": dtypes.FormatTypeVideoAudio,
		"video_mp4_h265": dtypes.FormatTypeVideoAudio,
		"audio_orig":     dtypes.FormatTypeAudioOnly,
		"audio_mp3":      dtypes.FormatTypeAudioOnly,
		"audio_m4a":      dtypes.FormatTypeAudioOnly,
	}
	videoFormatMap = map[string]dtypes.VideoFormat{
		"video_orig":     dtypes.VideoFormatOrig,
		"video_mp4_orig": dtypes.VideoFormatMP4Orig,
		"video_mp4_h264": dtypes.VideoFormatMP4H264,
		"video_mp4_h265": dtypes.VideoFormatMP4H265,
	}
	audioFormatMap = map[string]dtypes.AudioFormat{
		"audio_orig": dtypes.AudioFormatOrig,
		"audio_mp3":  dtypes.AudioFormatMP3,
		"audio_m4a":  dtypes.AudioFormatM4A,
	}
)

func (h *GrabberHandlers) GrabHandler(ctx *fasthttp.RequestCtx) {
	var itemsOnlyOne = false
	if itemsOnlyOneStr := string(ctx.Request.Header.Cookie("resultItemsOnlyOne")); itemsOnlyOneStr == "true" {
		itemsOnlyOne = true
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

	// Read selected format from radio buttons
	formFormatType := string(ctx.FormValue(formFieldFormatTypeKey))
	formatType := formatTypeMap[formFormatType]
	if formatType == "" {
		formatType = dtypes.FormatTypeVideoAudio
	}

	var videoFormat *dtypes.VideoFormat
	if _, exists := videoFormatMap[formFormatType]; exists {
		videoFormat = new(dtypes.VideoFormat)
		*videoFormat = videoFormatMap[formFormatType]
	}

	var audioFormat *dtypes.AudioFormat
	if _, exists := audioFormatMap[formFormatType]; exists {
		audioFormat = new(dtypes.AudioFormat)
		*audioFormat = audioFormatMap[formFormatType]
	}

	var (
		data = grabResultData{
			YoutubeURL: url,
		}
	)

	resp, err := h.usecases.Downloader.ScheduleDownload(
		ctx,
		url,
		&ddownload.DownloadOptions{
			FormatType:  formatType,
			VideoFormat: videoFormat,
			AudioFormat: audioFormat,
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

	dataMap := htmxvalues.MergeMaps(
		htmxvalues.PathValues,
		htmxvalues.IconNames,
		htmxvalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	dataMap[htmxvalues.GrabResultStatusIconNameKey] = htmxvalues.GrabResultStatusIconName(resp.Status)
	dataMap[htmxvalues.GrabResultItemStatusHtmlKey] = template.HTML(htmxvalues.GrabResultStatusIconSvgRaw(resp.Status, iconsDir))
	dataMap[htmxvalues.IsItemHTMXOptionRepeatKey] = true
	dataMap[htmxvalues.IsItemFirstKey] = itemsOnlyOne
	dataMap[htmxvalues.DataOnlyOneKey] = false
	dataMap[htmxvalues.ItemFadeKey] = "fade-in"
	dataMap[htmxvalues.GrabResultItemStatusTextKey] = ""
	dataMap[htmxvalues.ResultYoutubeUrlFadeKey] = ""
	dataMap[htmxvalues.ResultSizeFadeKey] = ""
	dataMap[htmxvalues.ResultFormatFadeKey] = ""

	var buf bytes.Buffer
	h.templates.ExecuteTemplate(&buf, htmxvalues.GrabResultNewItemHtmlFileName, dataMap)
	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
