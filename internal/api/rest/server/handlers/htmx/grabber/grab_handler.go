package grabberh

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/valyala/fasthttp"
)

const (
	formFieldYouTubeURL = "youtubeURL"
	formFieldFormatType = "formatType"
)

type resultData struct {
	YoutubeURL  string
	Format      string
	DownloadURL string
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
	url := string(ctx.FormValue(formFieldYouTubeURL))
	if url == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("URL is required")
		return
	}

	// Read selected format from radio buttons
	formFormatType := string(ctx.FormValue(formFieldFormatType))
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
		data     = resultData{YoutubeURL: url}
	)

	resp, err := h.usecases.Downloader.Download(
		url,
		&ddownload.DownloadOptions{
			FormatType:  formatType,
			VideoFormat: videoFormat,
			AudioFormat: audioFormat,
		},
	)

	if err != nil {
		tmplPath = filepath.Join(h.assetsDir, "templates", "grab_result_error.html")
	} else {
		tmplPath = filepath.Join(h.assetsDir, "templates", "grab_result_success.html")
		data.Format = resp.FileExt
		// Set URL for download endpoint
		data.DownloadURL = fmt.Sprintf("%s?file=%s", httppaths.GroupDownloader+httppaths.PathDownload, resp.FileFullName)
	}

	tpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Template error")
		return
	}

	var buf bytes.Buffer
	tpl.Execute(&buf, data)
	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
