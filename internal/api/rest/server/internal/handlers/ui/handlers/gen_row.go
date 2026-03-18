package handlers

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	dltypes "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/types"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/errorx"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

type fileRowInfoData struct {
	FileID           string
	PathFileRow      string
	PathFileRepeat   string
	YoutubeChannelID string
	AvatarTitle      string
	MediaTitle       string
	MediaURL         string
	FileSize         string
	Format           string
	DataFormat       string
	FormatTitle      string
	FormatTooltip    string
	IsAudio          string
	DownloadURL      string
	StreamURL        string
	DeleteURL        string
	LogoVersion      string
	RowID            string
	ProgressID       string
}

type genRowResult struct {
	templateName string
	data         map[string]any
	httpStatus   int
	err          error
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

	youtubeChannelID := channelIDValueNone
	if fileInfo.YoutubeChannelID != nil {
		youtubeChannelID = *fileInfo.YoutubeChannelID
	}

	logoVersion := ""
	if fileInfo.YoutubeChannelID != nil {
		logoVersion = "yt-channel"
	} else if fileInfo.HasSiteIcon {
		logoVersion = "site-logo"
	}
	if logoVersion == "" {
		logoVersion = fmt.Sprintf("%d", time.Now().UTC().Unix())
	}

	data := fileRowInfoData{
		FileID:           fileInfo.FileID.String(),
		MediaURL:         fileInfo.MediaUrl,
		YoutubeChannelID: youtubeChannelID,
		AvatarTitle:      fileInfo.AvatarTitle,
		MediaTitle:       fileInfo.MediaTitle,
		PathFileRow:      httppaths.BuildPathFileRow(fileInfo.FileID),
		PathFileRepeat:   httppaths.BuildPathFileRepeat(fileInfo.FileID),
		FileSize:         "-",
		Format:           "-",
		DataFormat:       "-",
		FormatTitle:      fileInfo.MediaInfoText,
		FormatTooltip:    fileInfo.MediaInfoTooltip,
		DownloadURL:      httppaths.BuildPathFileDownload(fileInfo.FileID),
		StreamURL:        httppaths.BuildPathFileStream(fileInfo.FileID),
		DeleteURL:        httppaths.BuildPathFile(fileInfo.FileID),
		LogoVersion:      logoVersion,
		RowID:            "row-" + fileInfo.FileID.String(),
		ProgressID:       "progress-" + fileInfo.FileID.String(),
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

	dataMap[uivalues.GrabResultStatusIconNameKey] = uivalues.DownloadResultStatusIconFileName(fileInfo.Status)
	dataMap[uivalues.IsItemHTMXOptionRepeatKey] = isGrabResultItemHTMXOptionRepeat
	dataMap[uivalues.IsFileEventKey] = isFileEvent
	dataMap[uivalues.GrabResultItemStatusHtmlKey] = template.HTML(
		uivalues.DownloadResultStatusIconSvgRaw(fileInfo.Status, iconsDir),
	)
	dataMap[uivalues.GrabResultItemStatusTextKey] = fileInfo.StatusText
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

	dataMap[uivalues.DownloadStatusKey] = fileInfo.Status.String()
	dataMap[uivalues.DownloadWorkingStatusKey] = dltypes.MapUsecaseWorkingStatusToUI(fileInfo.WorkingStatus)
	if fileInfo.Progress != nil {
		dataMap[uivalues.DownloadingProgressPercentKey] = int(fileInfo.Progress.Percent())
	}

	var tmplFileName = uivalues.GrabResultRowStatusHtmlFileName
	if fileInfo.Status == dtypes.FileStatusDone {
		tmplFileName = uivalues.GrabResultRowSuccessHtmlFileName
	}

	return genRowResult{
		templateName: tmplFileName,
		data:         dataMap,
		httpStatus:   fasthttp.StatusOK,
		err:          nil,
	}
}
