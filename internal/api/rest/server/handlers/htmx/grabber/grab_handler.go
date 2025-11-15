package grabberh

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	avalues "github.com/neosy/elengrab/internal/api/rest/server/assets/values"
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
		"video_orig": dtypes.FormatTypeVideoAudio,
		"video_mp4":  dtypes.FormatTypeVideoAudio,
		"audio_orig": dtypes.FormatTypeAudioOnly,
		"audio_mp3":  dtypes.FormatTypeAudioOnly,
		"audio_m4a":  dtypes.FormatTypeAudioOnly,
	}
	videoFormatMap = map[string]dtypes.VideoFormat{
		"video_orig": dtypes.VideoFormatOrig,
		"video_mp4":  dtypes.VideoFormatMP4,
	}
	audioFormatMap = map[string]dtypes.AudioFormat{
		"audio_orig": dtypes.AudioFormatOrig,
		"audio_mp3":  dtypes.AudioFormatMP3,
		"audio_m4a":  dtypes.AudioFormatM4A,
	}
)

func (h *GrabberHandlers) GrabHandler(ctx *fasthttp.RequestCtx) {
	itemsOnlyOne := string(ctx.Request.Header.Cookie("resultItemsOnlyOne"))

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
		tmplPath string
		data     = grabResultData{
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
		tmplPath = filepath.Join(h.assetsDir, "templates", "grab_result_failed.html")
	} else {
		if itemsOnlyOne == "true" {
			tmplPath = filepath.Join(h.assetsDir, "templates", "grab_result_first_new_item.html")
		} else {
			tmplPath = filepath.Join(h.assetsDir, "templates", "grab_result_new_item.html")
		}
		data.YoutubeTitle = url
		data.PathFileRow = httppaths.BuildPathFileRow(resp.FileId)
		data.FileSize = "-"
		data.Format = "-"
		// Set URL for download endpoint
		data.DownloadURL = fmt.Sprintf("%s?file=%s", httppaths.GroupDownloader+httppaths.PathDownload, resp.FileId)
	}

	tpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Template error")
		return
	}

	dataMap := avalues.MergeMaps(
		avalues.PathValues,
		avalues.IconNames,
		avalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	dataMap[avalues.GrabResultStatusIconNameKey] = avalues.GrabResultStatusIconName(resp.Status)
	dataMap[avalues.GrabResultItemStatusHtmlKey] = avalues.GrabResultStatusIconSvgRaw(resp.Status, iconsDir)

	var buf bytes.Buffer
	tpl.Execute(&buf, dataMap)
	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
