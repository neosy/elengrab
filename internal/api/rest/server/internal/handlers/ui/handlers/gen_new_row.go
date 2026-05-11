package handlers

import (
	"bytes"
	"html/template"
	"path/filepath"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type newRowData struct {
	FileID           string
	DownloadStatus   string
	RowID            string
	PathFileRow      string
	YoutubeChannelID string
	MediaTitle       string
	MediaURL         string
	FileSize         string
	Format           string
	ImageURL         string
	DeleteURL        string
}

func (h *DownloaderHandlers) genNewRow(
	fileInfo *ucdto.ScheduleDownloadResponse,
	pageHasDivItems bool,
) *bytes.Buffer {
	imageSources := []dtypes.ImageSource{
		dtypes.ImageSourceThumbnail,
		dtypes.ImageSourceAvatar,
		dtypes.ImageSourceSite,
	}

	fileImageURL := httppaths.BuildPathFileImage(fileInfo.FileID, fileInfo.ImageMetaHash(), imageSources)

	data := newRowData{
		MediaURL:       fileInfo.URL,
		DownloadStatus: fileInfo.Status.String(),
		DeleteURL:      httppaths.BuildPathFile(httppaths.PathFile, fileInfo.FileID),
		ImageURL:       fileImageURL,
		FileID:         fileInfo.FileID.String(),
		RowID:          "row-" + fileInfo.FileID.String(),
		MediaTitle:     fileInfo.URL,
		PathFileRow:    httppaths.BuildPathFileRow(fileInfo.FileID),
		FileSize:       "-",
		Format:         "-",
	}

	dataMap := uivalues.MergeMaps(
		uivalues.PathValues,
		uivalues.IconFileNames(),
		uivalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	dataMap[uivalues.ResultRowStatusIconKey] = template.HTML(
		uivalues.DownloaderResultStatusIconSvgRaw(fileInfo.Status, iconsDir))
	dataMap[uivalues.DownloaderResultItemDeleteIconKey] = template.HTML(
		uivalues.IconFileRaw(uivalues.IconFileName(uivalues.DownloadDeleteIconNameKey), iconsDir))
	dataMap[uivalues.IsItemHTMXOptionRepeatKey] = true
	dataMap[uivalues.PageHasDivItemsKey] = pageHasDivItems
	dataMap[uivalues.ResultRowFadeKey] = "fade-in"
	dataMap[uivalues.ResultRowStatusTitleKey] = ""
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
