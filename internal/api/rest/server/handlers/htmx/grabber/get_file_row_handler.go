package grabberh

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/google/uuid"
	avalues "github.com/neosy/elengrab/internal/api/rest/server/assets/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
	"github.com/valyala/fasthttp"
)

type fileRowInfoData struct {
	PathFileRow  string
	YoutubeTitle string
	YoutubeURL   string
	Format       string
	DownloadURL  string
}

func (h *GrabberHandlers) GetFileRow(ctx *fasthttp.RequestCtx) {
	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is required")
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is incorrect")
		return
	}

	var (
		tmplPath string
		data     fileRowInfoData
	)

	resp, err := h.usecases.Downloader.GetFileInfo(ctx, fileId)
	if err != nil {
		tmplPath = filepath.Join(h.assetsDir, "templates", "grab_result_failed.html")
	}

	if resp != nil {
		tmplPath = filepath.Join(h.assetsDir, "templates", avalues.GrabResultStatusHtmlFileName(resp.Status))

		data.YoutubeURL = resp.YoutubeUrl
		data.YoutubeTitle = resp.YoutubeTitle
		data.PathFileRow = httppaths.BuildPathFileRow(resp.FileId)
		data.Format = resp.FileExt
		if data.Format == "" {
			data.Format = "-"
		}
		// Set URL for download endpoint
		data.DownloadURL = httppaths.BuildPathFileDownload(resp.FileId)
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
	dataMap[avalues.GrabResultItemHtmxOptionKey] = avalues.GrabResultStatusHtmxOption(resp.Status, data.PathFileRow)
	dataMap[avalues.GrabResultItemStatusHtmlKey] = avalues.GrabResultStatusIconSvgRaw(resp.Status, iconsDir)
	dataMap[avalues.GrabResultItemStatusTextKey] = resp.StatusText

	var buf bytes.Buffer
	tpl.Execute(&buf, dataMap)
	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
