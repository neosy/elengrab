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
	FileID         string
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
	fileInfo *dto.GetFileInfoResponse,
	isFileEvent bool,
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

	if fileInfo == nil {
		return genRowResult{
			httpStatus: fasthttp.StatusInternalServerError,
			err:        errorx.New("the request returned an empty"),
		}
	}

	var youtubeChannelID string
	if fileInfo.ChannelID != nil && fileInfo.IsYouTube() {
		youtubeChannelID = *fileInfo.ChannelID
	}

	var thumbnailID string
	if fileInfo.MediaInfo != nil && fileInfo.MediaInfo.GetThumbnailID() != nil {
		thumbnailID = fileInfo.MediaInfo.GetThumbnailID().String()
	}

	fileImageURL := httppaths.BuildPathFileImage(
		fileInfo.FileID,
		fileInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceThumbnail,
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		},
	)

	fileImageAvatarURL := httppaths.BuildPathFileImage(
		fileInfo.FileID,
		fileInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		},
	)

	fileImageSiteURL := httppaths.BuildPathFileImage(
		fileInfo.FileID,
		fileInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceSite,
		},
	)

	data := fileRowInfoData{
		FileID:         fileInfo.FileID.String(),
		DownloadStatus: fileInfo.Status.String(),
		WorkingStatus:  dltypes.MapUsecaseWorkingStatusToUI(fileInfo.WorkingStatus).String(),

		YoutubeChannelID: youtubeChannelID,
		ThumbnailID:      thumbnailID,
		AvatarTitle:      fileInfo.AvatarTitle,

		MediaTitle: fileInfo.MediaTitle,
		MediaURL:   fileInfo.MediaURL,

		ContentTimeAgo: fileInfo.CreatedTimeAgo,

		ImageURL:       fileImageURL,
		ImageAvatarURL: fileImageAvatarURL,
		ImageSiteURL:   fileImageSiteURL,

		PathFileRow:    httppaths.BuildPathFileRow(fileInfo.FileID),
		PathFileRepeat: httppaths.BuildPathFileRepeat(fileInfo.FileID),

		FileSize:      "-",
		Format:        "-",
		DataFormat:    "-",
		FormatTitle:   fileInfo.MediaInfoText,
		FormatTooltip: fileInfo.MediaInfoTooltip,

		DownloadURL: httppaths.BuildPathFileDownload(fileInfo.FileID),
		WatchURL:    httppaths.BuildPathFileWatch(fileInfo.FileID),
		StreamURL:   httppaths.BuildPathFileStream(fileInfo.FileID),
		DeleteURL:   httppaths.BuildPathFile(httppaths.PathFile, fileInfo.FileID),

		RowID:      "row-" + fileInfo.FileID.String(),
		ProgressID: "progress-" + fileInfo.FileID.String(),
	}

	if fileInfo.FileSize != nil && *fileInfo.FileSize > 0 {
		data.FileSize = uformat.BytesHuman(*fileInfo.FileSize)
	}
	if fileInfo.FileExt != "" {
		data.Format = fileInfo.FileExt
		data.DataFormat = fileInfo.FileExt
	}
	if fileInfo.MediaInfo != nil {
		data.IsAudio = fmt.Sprint(fileInfo.MediaInfo.FormatType == dtypes.FormatTypeAudioOnly)
	}

	dataMap := uivalues.MergeMaps(
		uivalues.PathValues,
		uivalues.IconFileNames(),
		uivalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	var isGrabResultItemHTMXOptionRepeat = false
	switch fileInfo.Status {
	case dtypes.FileStatusNew, dtypes.FileStatusPending, dtypes.FileStatusWorking:
		isGrabResultItemHTMXOptionRepeat = true
	}

	dataMap[uivalues.IsItemHTMXOptionRepeatKey] = isGrabResultItemHTMXOptionRepeat
	dataMap[uivalues.IsFileEventKey] = isFileEvent
	dataMap[uivalues.ResultRowStatusIconKey] = template.HTML(
		uivalues.DownloadResultStatusIconSvgRaw(fileInfo.Status, iconsDir),
	)
	dataMap[uivalues.ResultRowStatusTitleKey] = fileInfo.StatusText
	dataMap[uivalues.DownloadResultItemDeleteIconKey] = template.HTML(
		uivalues.IconFileRaw(uivalues.IconFileName(uivalues.DownloadDeleteIconNameKey), iconsDir),
	)

	if fileInfo.Status == dtypes.FileStatusFailed {
		dataMap[uivalues.DownloadResultItemStatusFailedIconKey] = template.HTML(
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
	if fileInfo.Status == dtypes.FileStatusWorking {
		dataMap[uivalues.IsItemSpinerKey] = true
	}

	if fileInfo.Progress != nil {
		dataMap[uivalues.DownloadingProgressPercentKey] = int(fileInfo.Progress.Percent())
	}

	return genRowResult{
		data:       dataMap,
		httpStatus: fasthttp.StatusOK,
		err:        nil,
	}
}
