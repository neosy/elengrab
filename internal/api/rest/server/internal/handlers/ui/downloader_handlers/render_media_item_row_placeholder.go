package handlers

import (
	"bytes"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (h *DownloaderHandlers) renderMediaItemRowPlaceholder(
	downloadInfo *ucdto.ScheduleDownloadResponse,
	pageHasDivItems bool,
) *bytes.Buffer {
	imageSources := []dtypes.ImageSource{
		dtypes.ImageSourceThumbnail,
		dtypes.ImageSourceAvatar,
		dtypes.ImageSourceSite,
	}

	downloadImageURL := httppaths.BuildPathMediaItemImage(downloadInfo.DownloadID, downloadInfo.ImageMetaHash(time.Now().String()), imageSources)

	data := pages.RowFragmentValues{
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

		DownloaderResultItemStatusIcon: icons.DownloaderIconByStatus(downloadInfo.Status).FileRaw(),
		DownloaderResultItemDeleteIcon: icons.DownloadDeleteIcon.FileRaw(),
		IsItemHTMXOptionRepeat:         true,
		PageHasDivItems:                pageHasDivItems,
		ResultRowFade:                  "fade-in",
		ResultRowStatusTitle:           "",
		ResultMediaUrlFade:             "",
		ResultSizeFade:                 "",
		ResultFormatFade:               "",
	}

	pageData := pages.RowFragmentData{
		BasePaths:     paths.NewPaths(),
		Values:        &data,
		IconFileNames: icons.FileNamesByKey(),
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, components.ResultNewRowKey, pageData)
	if err != nil {
		return nil
	}

	return &buf
}
