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
	"github.com/neosy/elengrab/internal/pkg/humanize"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

func (h *DownloaderHandlers) renderMediaItemRowPlaceholder(
	schelduleDownload *ucdto.ScheduleDownloadResponse,
	pageHasDivItems bool,
) *bytes.Buffer {
	imageSources := []dtypes.ImageSource{
		dtypes.ImageSourceThumbnail,
		dtypes.ImageSourceAvatar,
		dtypes.ImageSourceSite,
	}

	downloadImageURL := httppaths.BuildMediaItemImagePath(schelduleDownload.DownloadID, schelduleDownload.ImageMetaHash(time.Now().String()), imageSources)

	id := idcodec.EncodeUUIDBase64URL(schelduleDownload.DownloadID)

	data := pages.RowFragmentValues{
		DownloadID:     id,
		MediaURL:       schelduleDownload.URL,
		DownloadStatus: schelduleDownload.Status.String(),
		LazyLoadImages: false,

		RowID:      "row-" + id,
		MediaTitle: schelduleDownload.URL,

		ThumbnailURL: h.thumbnailURLWithFallback(nil),

		ImageURL: downloadImageURL,

		ImageAvatarURL: httppaths.BuildMediaItemImagePath(
			schelduleDownload.DownloadID,
			schelduleDownload.ImageMetaHash(time.Now().String()),
			[]dtypes.ImageSource{dtypes.ImageSourceAvatar},
		),

		DeleteURL: httppaths.BuildMediaItemPath(schelduleDownload.DownloadID),
		FilePath:  httppaths.BuildMediaItemRowPath(schelduleDownload.DownloadID),
		FileSize:  "-",
		Format:    "-",
		Duration:  humanize.DurationClock(0),

		DownloaderResultItemStatusIcon: icons.DownloaderIconByStatus(schelduleDownload.Status).FileRaw(),
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
