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

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	data := uivalues.RowFragmentValues{
		MediaURL:       downloadInfo.URL,
		DownloadStatus: downloadInfo.Status.String(),
		DeleteURL:      httppaths.BuildPathMediaItem(httppaths.PathMediaItem, downloadInfo.DownloadID),
		ImageURL:       downloadImageURL,
		DownloadID:     downloadInfo.DownloadID.String(),
		RowID:          "row-" + downloadInfo.DownloadID.String(),
		MediaTitle:     downloadInfo.URL,
		FilePath:       httppaths.BuildPathMediaItemRow(downloadInfo.DownloadID),
		FileSize:       "-",
		Format:         "-",

		DownloaderResultItemStatusIcon: template.HTML(uivalues.DownloaderResultStatusIconSvgRaw(downloadInfo.Status, iconsDir)),
		DownloaderResultItemDeleteIcon: template.HTML(uivalues.IconFileRaw(uivalues.IconFileName(uivalues.DownloadDeleteIconNameKey), iconsDir)),
		IsItemHTMXOptionRepeat:         true,
		PageHasDivItems:                pageHasDivItems,
		ResultRowFade:                  "fade-in",
		ResultRowStatusTitle:           "",
		ResultMediaUrlFade:             "",
		ResultSizeFade:                 "",
		ResultFormatFade:               "",
	}

	pageData := uivalues.RowFragmentData{
		BasePaths:     uivalues.NewBasePaths(),
		Values:        data,
		IconFileNames: uivalues.IconFileNames(),
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, uivalues.ComponentResultNewRowKey, pageData)
	if err != nil {
		return nil
	}

	return &buf
}
