package handlers

import (
	"bytes"
	"html/template"
	"path/filepath"
	"time"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type newRowData struct {
	DownloadID       string
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
	downloadInfo *ucdto.ScheduleDownloadResponse,
	pageHasDivItems bool,
) *bytes.Buffer {
	imageSources := []dtypes.ImageSource{
		dtypes.ImageSourceThumbnail,
		dtypes.ImageSourceAvatar,
		dtypes.ImageSourceSite,
	}

	downloadImageURL := httppaths.BuildPathMediaItemImage(downloadInfo.DownloadID, downloadInfo.ImageMetaHash(time.Now().String()), imageSources)

	data := newRowData{
		MediaURL:       downloadInfo.URL,
		DownloadStatus: downloadInfo.Status.String(),
		DeleteURL:      httppaths.BuildPathMediaItem(httppaths.PathMediaItem, downloadInfo.DownloadID),
		ImageURL:       downloadImageURL,
		DownloadID:     downloadInfo.DownloadID.String(),
		RowID:          "row-" + downloadInfo.DownloadID.String(),
		MediaTitle:     downloadInfo.URL,
		PathFileRow:    httppaths.BuildPathMediaItemRow(downloadInfo.DownloadID),
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
		uivalues.DownloaderResultStatusIconSvgRaw(downloadInfo.Status, iconsDir))
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
