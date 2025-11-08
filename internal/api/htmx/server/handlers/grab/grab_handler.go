package grabh

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

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
	mapFormatType = map[string]dtypes.FormatType{
		"video": dtypes.FormatTypeVideoAudio,
		"audio": dtypes.FormatTypeAudioOnly,
	}
)

func (h *GrabHandlers) GrabHandler(ctx *fasthttp.RequestCtx) {
	url := string(ctx.FormValue(formFieldYouTubeURL))
	if url == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("URL is required")
		return
	}

	// Read selected format from radio buttons
	formatType := mapFormatType[string(ctx.FormValue(formFieldFormatType))]
	if formatType == "" {
		formatType = dtypes.FormatTypeVideoAudio
	}

	var (
		tmplPath string
		data     = resultData{YoutubeURL: url}
	)

	resp, err := h.usecases.Downloader.Download(url, &ddownload.DownloadOptions{FormatType: formatType})

	if err != nil {
		tmplPath = filepath.Join(h.assetsDir, "templates", "grab_result_error.html")
	} else {
		tmplPath = filepath.Join(h.assetsDir, "templates", "grab_result_success.html")
		data.Format = resp.FileExt
		// Set URL for download endpoint
		data.DownloadURL = fmt.Sprintf("/downloader/download?file=%s", resp.FileFullName)
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
