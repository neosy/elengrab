package handlers

import (
	"fmt"
	"html/template"
	"path/filepath"

	dltypes "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/types"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

type fileRowInfoData struct {
	DownloadID     string
	DownloadStatus string
	WorkingStatus  string

	PathFileRow    string
	PathFileRepeat string

	YoutubeChannelID string
	ThumbnailID      string
	AvatarTitle      string

	ImageURL       string
	ImageAvatarURL string
	ImageSiteURL   string

	MediaTitle string
	MediaURL   string

	ContentTimeAgo string

	FileSize      string
	Format        string
	DataFormat    string
	FormatTitle   string
	FormatTooltip string
	IsAudio       string
	DownloadURL   string
	StreamURL     string
	WatchURL      string
	DeleteURL     string
	RowID         string
	ProgressID    string
}

type genRowResult struct {
	data       map[string]any
	httpStatus int
	err        error
}

func (h *DownloaderHandlers) genRow(
	downloadInfo *dto.GetMediaDownloadInfoResponse,
	isDownloadEvent bool,
) genRowResult {
	var (
		cacheChanged = struct {
			youtubeChannelID bool
			mediaTitle       bool
			FileSize         bool
			Format           bool
			Status           bool
			ProgressPercent  bool
		}{}
	)

	if downloadInfo == nil {
		return genRowResult{
			httpStatus: fasthttp.StatusInternalServerError,
			err:        errorx.New("the request returned an empty"),
		}
	}

	var youtubeChannelID string
	if downloadInfo.ChannelID != nil && downloadInfo.IsYouTube() {
		youtubeChannelID = *downloadInfo.ChannelID
	}

	var thumbnailID string
	if downloadInfo.MediaInfo != nil && downloadInfo.MediaInfo.GetThumbnailID() != nil {
		thumbnailID = downloadInfo.MediaInfo.GetThumbnailID().String()
	}

	fileImageURL := httppaths.BuildPathFileImage(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceThumbnail,
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		},
	)

	fileImageAvatarURL := httppaths.BuildPathFileImage(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		},
	)

	fileImageSiteURL := httppaths.BuildPathFileImage(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceSite,
		},
	)

	data := fileRowInfoData{
		DownloadID:     downloadInfo.DownloadID.String(),
		DownloadStatus: downloadInfo.Status.String(),
		WorkingStatus:  dltypes.MapUsecaseWorkingStatusToUI(downloadInfo.WorkingStatus).String(),

		YoutubeChannelID: youtubeChannelID,
		ThumbnailID:      thumbnailID,
		AvatarTitle:      downloadInfo.AvatarTitle,

		MediaTitle: downloadInfo.MediaTitle,
		MediaURL:   downloadInfo.MediaURL,

		ContentTimeAgo: downloadInfo.CreatedTimeAgo,

		ImageURL:       fileImageURL,
		ImageAvatarURL: fileImageAvatarURL,
		ImageSiteURL:   fileImageSiteURL,

		PathFileRow:    httppaths.BuildPathFileRow(downloadInfo.DownloadID),
		PathFileRepeat: httppaths.BuildPathFileRepeat(downloadInfo.DownloadID),

		FileSize:      "-",
		Format:        "-",
		DataFormat:    "-",
		FormatTitle:   downloadInfo.MediaInfoText,
		FormatTooltip: downloadInfo.MediaInfoTooltip,

		DownloadURL: httppaths.BuildPathFileDownload(downloadInfo.DownloadID),
		WatchURL:    httppaths.BuildPathFileWatch(downloadInfo.DownloadID),
		StreamURL:   httppaths.BuildPathFileStream(downloadInfo.DownloadID),
		DeleteURL:   httppaths.BuildPathFile(httppaths.PathFile, downloadInfo.DownloadID),

		RowID:      "row-" + downloadInfo.DownloadID.String(),
		ProgressID: "progress-" + downloadInfo.DownloadID.String(),
	}

	if downloadInfo.FileSize != nil && *downloadInfo.FileSize > 0 {
		data.FileSize = uformat.BytesHuman(*downloadInfo.FileSize)
	}
	if downloadInfo.FileExt != "" {
		data.Format = downloadInfo.FileExt
		data.DataFormat = downloadInfo.FileExt
	}
	if downloadInfo.MediaInfo != nil {
		data.IsAudio = fmt.Sprint(downloadInfo.MediaInfo.FormatType == dtypes.FormatTypeAudioOnly)
	}

	dataMap := uivalues.MergeMaps(
		uivalues.PathValues,
		uivalues.IconFileNames(),
		uivalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	var isGrabResultItemHTMXOptionRepeat = false
	switch downloadInfo.Status {
	case dtypes.MediaDownloadStatusNew, dtypes.MediaDownloadStatusPending, dtypes.MediaDownloadStatusWorking:
		isGrabResultItemHTMXOptionRepeat = true
	}

	dataMap[uivalues.IsItemHTMXOptionRepeatKey] = isGrabResultItemHTMXOptionRepeat
	dataMap[uivalues.IsFileEventKey] = isDownloadEvent
	dataMap[uivalues.ResultRowStatusIconKey] = template.HTML(
		uivalues.DownloaderResultStatusIconSvgRaw(downloadInfo.Status, iconsDir),
	)
	dataMap[uivalues.ResultRowStatusTitleKey] = downloadInfo.StatusText
	dataMap[uivalues.DownloaderResultItemDeleteIconKey] = template.HTML(
		uivalues.IconFileRaw(uivalues.IconFileName(uivalues.DownloadDeleteIconNameKey), iconsDir),
	)

	if downloadInfo.Status == dtypes.MediaDownloadStatusFailed {
		dataMap[uivalues.DownloaderResultItemStatusFailedIconKey] = template.HTML(
			uivalues.IconFileRaw(uivalues.IconFileName(uivalues.DownloadRepeatIconNameKey), iconsDir),
		)
	}

	dataMap[uivalues.ResultMediaUrlFadeKey] = ""
	dataMap[uivalues.ResultSizeFadeKey] = ""
	dataMap[uivalues.ResultFormatFadeKey] = ""
	if cacheChanged.mediaTitle {
		dataMap[uivalues.ResultMediaUrlFadeKey] = "fade-text"
	}
	if cacheChanged.FileSize {
		dataMap[uivalues.ResultSizeFadeKey] = "fade-text"
	}
	if cacheChanged.Format {
		dataMap[uivalues.ResultFormatFadeKey] = "fade-text"
	}

	dataMap[uivalues.IsItemSpinerKey] = false
	if downloadInfo.Status == dtypes.MediaDownloadStatusWorking {
		dataMap[uivalues.IsItemSpinerKey] = true
	}

	if downloadInfo.Progress != nil {
		dataMap[uivalues.DownloadingProgressPercentKey] = int(downloadInfo.Progress.Percent())
	}

	return genRowResult{
		data:       dataMap,
		httpStatus: fasthttp.StatusOK,
		err:        nil,
	}
}
