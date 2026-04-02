package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
)

type newRowData struct {
	FileID           string
	RowID            string
	PathFileRow      string
	YoutubeChannelID string
	MediaTitle       string
	MediaURL         string
	FileSize         string
	Format           string
	LogoVersion      string
	DeleteURL        string
}

func (h *DownloaderHandlers) genNewRow(
	fileInfo *ucdto.ScheduleDownloadResponse,
	pageHasDivItems bool,
) *bytes.Buffer {
	data := newRowData{
		MediaURL:         fileInfo.URL,
		YoutubeChannelID: channelIDValueNone,
		LogoVersion:      fmt.Sprintf("%d", time.Now().UTC().Unix()),
		DeleteURL:        httppaths.BuildPathFile(fileInfo.FileID),
		FileID:           fileInfo.FileID.String(),
		RowID:            "row-" + fileInfo.FileID.String(),
		MediaTitle:       fileInfo.URL,
		PathFileRow:      httppaths.BuildPathFileRow(fileInfo.FileID),
		FileSize:         "-",
		Format:           "-",
	}

	dataMap := uivalues.MergeMaps(
		uivalues.PathValues,
		uivalues.IconFileNames(),
		uivalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	dataMap[uivalues.GrabResultStatusIconNameKey] = uivalues.DownloadResultStatusIconFileName(fileInfo.Status)
	dataMap[uivalues.GrabResultItemStatusHtmlKey] = template.HTML(
		uivalues.DownloadResultStatusIconSvgRaw(fileInfo.Status, iconsDir))
	dataMap[uivalues.DownloadResultItemDeleteIconKey] = template.HTML(
		uivalues.IconFileRaw(uivalues.IconFileName(uivalues.DownloadDeleteIconNameKey), iconsDir))
	dataMap[uivalues.IsItemHTMXOptionRepeatKey] = true
	dataMap[uivalues.PageHasDivItemsKey] = pageHasDivItems
	dataMap[uivalues.ItemFadeKey] = "fade-in"
	dataMap[uivalues.GrabResultItemStatusTextKey] = ""
	dataMap[uivalues.ResultMediaUrlFadeKey] = ""
	dataMap[uivalues.ResultSizeFadeKey] = ""
	dataMap[uivalues.ResultFormatFadeKey] = ""

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, uivalues.ComponentResultNewRowKey, dataMap)
	if err != nil {
		return nil
	}

	return &buf
}
