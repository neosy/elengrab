package grabberh

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	htmxvalues "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/utils"
	"github.com/valyala/fasthttp"
)

type fileRowInfoData struct {
	PathFileRow  string
	YoutubeTitle string
	YoutubeURL   string
	FileSize     string
	Format       string
	DataFormat   string
	FormatTitle  string
	DownloadURL  string
}

type cacheRowEntry struct {
	youtubeTitle string
	FileSize     int
	Format       string
	Status       dtypes.FileStatus
	Updated      time.Time
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

func (h *GrabberHandlers) genRow(ctx context.Context, fileInfo *dto.GetFileInfoResponse, isLoadHistory bool) (*bytes.Buffer, int, error) {
	var (
		cacheChanged = struct {
			youtubeTitle bool
			FileSize     bool
			Format       bool
			Status       bool
		}{}
	)

	if fileInfo == nil {
		return nil, fasthttp.StatusInternalServerError, fmt.Errorf("the request returned an empty")
	}

	// Checking the cache
	{
		cacheRow.mu.RLock()
		cached, exists := cacheRow.data[fileInfo.FileId]
		cacheRow.mu.RUnlock()

		if exists {
			cacheChanged.youtubeTitle = cached.youtubeTitle != fileInfo.YoutubeTitle
			if fileInfo.FileSize != nil {
				cacheChanged.FileSize = cached.FileSize != *fileInfo.FileSize
			}
			cacheChanged.Format = cached.Format != fileInfo.FileExt
			cacheChanged.Status = cached.Status != fileInfo.Status
		}

		if exists && fileInfo.Status != dtypes.FileStatusDone &&
			!cacheChanged.youtubeTitle &&
			!cacheChanged.FileSize &&
			!cacheChanged.Format &&
			!cacheChanged.Status &&
			time.Since(cached.Updated) < cacheRow.ttl {
			if !isLoadHistory {
				return nil, fasthttp.StatusNoContent, nil
			}
		}

		// Updating the cache
		var fileSize int
		if fileInfo.FileSize != nil {
			fileSize = *fileInfo.FileSize
		}
		cacheRow.mu.Lock()
		cacheRow.data[fileInfo.FileId] = cacheRowEntry{
			youtubeTitle: fileInfo.YoutubeTitle,
			FileSize:     fileSize,
			Format:       fileInfo.FileExt,
			Status:       fileInfo.Status,
			Updated:      time.Now(),
		}
		cacheRow.mu.Unlock()
	}

	data := fileRowInfoData{
		YoutubeURL:   fileInfo.YoutubeUrl,
		YoutubeTitle: fileInfo.YoutubeTitle,
		PathFileRow:  httppaths.BuildPathFileRow(fileInfo.FileId),
		FileSize:     "-",
		Format:       "-",
		DataFormat:   "-",
		FormatTitle:  fileInfo.MediaInfoText,
		DownloadURL:  httppaths.BuildPathFileDownload(fileInfo.FileId),
	}

	if fileInfo.FileSize != nil && *fileInfo.FileSize > 0 {
		data.FileSize = utils.BytesToHuman(*fileInfo.FileSize)
	}
	if fileInfo.FileExt != "" {
		data.Format = fileInfo.FileExt
		data.DataFormat = fileInfo.FileExt
	}

	dataMap := htmxvalues.MergeMaps(
		htmxvalues.PathValues,
		htmxvalues.IconNames,
		htmxvalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	var isGrabResultItemHTMXOptionRepeat = false
	switch fileInfo.Status {
	case dtypes.FileStatusNew, dtypes.FileStatusPending, dtypes.FileStatusWorking:
		isGrabResultItemHTMXOptionRepeat = true
	}

	dataMap[htmxvalues.GrabResultStatusIconNameKey] = htmxvalues.GrabResultStatusIconName(fileInfo.Status)
	dataMap[htmxvalues.IsItemHTMXOptionRepeatKey] = isGrabResultItemHTMXOptionRepeat
	dataMap[htmxvalues.GrabResultItemStatusHtmlKey] = template.HTML(htmxvalues.GrabResultStatusIconSvgRaw(fileInfo.Status, iconsDir))
	dataMap[htmxvalues.GrabResultItemStatusTextKey] = fileInfo.StatusText

	dataMap[htmxvalues.ResultYoutubeUrlFadeKey] = ""
	dataMap[htmxvalues.ResultSizeFadeKey] = ""
	dataMap[htmxvalues.ResultFormatFadeKey] = ""
	if cacheChanged.youtubeTitle {
		dataMap[htmxvalues.ResultYoutubeUrlFadeKey] = "fade-text"
	}
	if cacheChanged.FileSize {
		dataMap[htmxvalues.ResultSizeFadeKey] = "fade-text"
	}
	if cacheChanged.Format {
		dataMap[htmxvalues.ResultFormatFadeKey] = "fade-text"
	}

	dataMap[htmxvalues.IsItemSpinerKey] = false
	if fileInfo.Status == dtypes.FileStatusWorking {
		dataMap[htmxvalues.IsItemSpinerKey] = true
	}

	var tmplFileName = htmxvalues.GrabResultItemStatusHtmlFileName
	if fileInfo.Status == dtypes.FileStatusDone {
		tmplFileName = htmxvalues.GrabResultItemSuccessHtmlFileName
	}

	var buf bytes.Buffer
	err := h.templates.ExecuteTemplate(&buf, tmplFileName, dataMap)
	if err != nil {
		return nil, fasthttp.StatusInternalServerError, err
	}

	return &buf, fasthttp.StatusOK, nil
}
