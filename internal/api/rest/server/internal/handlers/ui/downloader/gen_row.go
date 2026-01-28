package downloaderh

import (
	"bytes"
	"html/template"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	dltypes "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/utils"
	"github.com/valyala/fasthttp"
)

type fileRowInfoData struct {
	PathFileRow      string
	PathFileRepeat   string
	YoutubeChannelID string
	YoutubeTitle     string
	YoutubeURL       string
	FileSize         string
	Format           string
	DataFormat       string
	FormatTitle      string
	DownloadURL      string
	DeleteURL        string
}

type cacheRowEntry struct {
	youtubeChannelID string
	youtubeTitle     string
	FileSize         int
	Format           string
	Status           dtypes.FileStatus
	ProgressPercent  int
	Updated          time.Time
}

type cache[T any] struct {
	data map[uuid.UUID]T
	mu   sync.RWMutex
	ttl  time.Duration
}

var (
	cacheRow = cache[cacheRowEntry]{
		data: make(map[uuid.UUID]cacheRowEntry),
		ttl:  time.Minute,
	}
)

func (h *DownloaderHandlers) genRow(fileInfo *dto.GetFileInfoResponse, isLoadHistory bool) (*bytes.Buffer, int, error) {
	var (
		cacheChanged = struct {
			youtubeChannelID bool
			youtubeTitle     bool
			FileSize         bool
			Format           bool
			Status           bool
			ProgressPercent  bool
		}{}
	)

	if fileInfo == nil {
		return nil, fasthttp.StatusInternalServerError, errorx.New("the request returned an empty")
	}

	// Checking the cache
	{
		cacheRow.mu.RLock()
		cached, exists := cacheRow.data[fileInfo.FileId]
		cacheRow.mu.RUnlock()

		if exists {
			if fileInfo.YoutubeChannelID != nil {
				cacheChanged.youtubeChannelID = cached.youtubeChannelID != *fileInfo.YoutubeChannelID
			}
			cacheChanged.youtubeTitle = cached.youtubeTitle != fileInfo.YoutubeTitle
			if fileInfo.FileSize != nil {
				cacheChanged.FileSize = cached.FileSize != *fileInfo.FileSize
			}
			cacheChanged.Format = cached.Format != fileInfo.FileExt
			cacheChanged.Status = cached.Status != fileInfo.Status
			if fileInfo.Progress != nil {
				cacheChanged.ProgressPercent = cached.ProgressPercent != int(fileInfo.Progress.Percent())
			}
		}

		if exists && fileInfo.Status != dtypes.FileStatusDone &&
			!cacheChanged.youtubeChannelID &&
			!cacheChanged.youtubeTitle &&
			!cacheChanged.FileSize &&
			!cacheChanged.Format &&
			!cacheChanged.Status &&
			!cacheChanged.ProgressPercent &&
			time.Since(cached.Updated) < cacheRow.ttl {
			if !isLoadHistory {
				return nil, fasthttp.StatusNoContent, nil
			}
		}

		// Updating the cache
		var youtubeChannelID string
		if fileInfo.YoutubeChannelID != nil {
			youtubeChannelID = *fileInfo.YoutubeChannelID
		}
		var fileSize int
		if fileInfo.FileSize != nil {
			fileSize = *fileInfo.FileSize
		}
		var progressPercent int
		if fileInfo.Progress != nil {
			progressPercent = int(fileInfo.Progress.Percent())
		}
		cacheRow.mu.Lock()
		cacheRow.data[fileInfo.FileId] = cacheRowEntry{
			youtubeChannelID: youtubeChannelID,
			youtubeTitle:     fileInfo.YoutubeTitle,
			FileSize:         fileSize,
			Format:           fileInfo.FileExt,
			Status:           fileInfo.Status,
			ProgressPercent:  progressPercent,
			Updated:          time.Now(),
		}
		cacheRow.mu.Unlock()
	}

	youtubeChannelID := channelIDValueNone
	if fileInfo.YoutubeChannelID != nil {
		youtubeChannelID = *fileInfo.YoutubeChannelID
	}

	data := fileRowInfoData{
		YoutubeURL:       fileInfo.YoutubeUrl,
		YoutubeChannelID: youtubeChannelID,
		YoutubeTitle:     fileInfo.YoutubeTitle,
		PathFileRow:      httppaths.BuildPathFileRow(fileInfo.FileId),
		PathFileRepeat:   httppaths.BuildPathFileRepeat(fileInfo.FileId),
		FileSize:         "-",
		Format:           "-",
		DataFormat:       "-",
		FormatTitle:      fileInfo.MediaInfoText,
		DownloadURL:      httppaths.BuildPathFileDownload(fileInfo.FileId),
		DeleteURL:        httppaths.BuildPathFileRow(fileInfo.FileId),
	}

	if fileInfo.FileSize != nil && *fileInfo.FileSize > 0 {
		data.FileSize = utils.BytesToHuman(*fileInfo.FileSize)
	}
	if fileInfo.FileExt != "" {
		data.Format = fileInfo.FileExt
		data.DataFormat = fileInfo.FileExt
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

	dataMap[uivalues.ResultYoutubeUrlFadeKey] = ""
	dataMap[uivalues.ResultSizeFadeKey] = ""
	dataMap[uivalues.ResultFormatFadeKey] = ""
	if cacheChanged.youtubeTitle {
		dataMap[uivalues.ResultYoutubeUrlFadeKey] = "fade-text"
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
		dataMap[uivalues.DownloadingProgressKey] = int(fileInfo.Progress.Percent())
	}

	var tmplFileName = uivalues.GrabResultItemStatusHtmlFileName
	if fileInfo.Status == dtypes.FileStatusDone {
		tmplFileName = uivalues.GrabResultItemSuccessHtmlFileName
	}

	var buf bytes.Buffer
	err := h.templates.ExecuteTemplate(&buf, tmplFileName, dataMap)
	if err != nil {
		return nil, fasthttp.StatusInternalServerError, errorx.NewByErr(err)
	}

	return &buf, fasthttp.StatusOK, nil
}
