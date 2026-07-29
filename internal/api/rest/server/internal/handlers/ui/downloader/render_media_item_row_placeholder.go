package downloader

import (
	"bytes"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
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

	id := idcodec.EncodeUUIDBase64URL(downloadInfo.DownloadID)

	data := pages.RowFragmentValues{
		MediaURL:       downloadInfo.URL,
		DownloadStatus: downloadInfo.Status.String(),
		DeleteURL:      httppaths.BuildMediaItemPath(downloadInfo.DownloadID),
		ImageURL:       downloadImageURL,
		DownloadID:     id,
		RowID:          "row-" + id,
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
		BasePaths:     paths.NewHttpPaths(),
		Values:        &data,
		IconFileNames: icons.FileNamesByKey(),
	}

	var buf bytes.Buffer
	err := h.templates.Base.ExecuteTemplate(&buf, components.ResultNewRowKey, pageData)
	if err != nil {
		return nil
	}

	return &buf
}
