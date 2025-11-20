package grabberh

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/google/uuid"
	avalues "github.com/neosy/elengrab/internal/api/rest/server/assets/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/utils"
	"github.com/valyala/fasthttp"
)

type cacheRowEntry struct {
	youtubeTitle string
	FileSize     int
	Format       string
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

func (h *GrabberHandlers) genFileRow(ctx context.Context, fileInfo *dto.GetFileInfoResponse, isLoadHistory bool) (*bytes.Buffer, int, error) {
	var (
		tmplPath string
		data     fileRowInfoData

		cacheChanged = struct {
			youtubeTitle bool
			FileSize     bool
			Format       bool
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
		}

		if exists && fileInfo.Status != dtypes.FileStatusDone &&
			!cacheChanged.youtubeTitle &&
			!cacheChanged.FileSize &&
			!cacheChanged.Format &&
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
			Updated:      time.Now(),
		}
		cacheRow.mu.Unlock()
	}

	tmplPath = filepath.Join(h.assetsDir, "templates", avalues.GetGrabResultStatusHtmlFileName(fileInfo.Status))

	data.YoutubeURL = fileInfo.YoutubeUrl
	data.YoutubeTitle = fileInfo.YoutubeTitle
	data.PathFileRow = httppaths.BuildPathFileRow(fileInfo.FileId)

	data.FileSize = "-"
	if fileInfo.FileSize != nil && *fileInfo.FileSize > 0 {
		data.FileSize = utils.BytesToHuman(*fileInfo.FileSize)
	}

	data.Format = "-"
	if fileInfo.FileExt != "" {
		data.Format = fileInfo.FileExt
	}
	// Set URL for download endpoint
	data.DownloadURL = httppaths.BuildPathFileDownload(fileInfo.FileId)

	tpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, fasthttp.StatusInternalServerError, fmt.Errorf("template error: %v", err)
	}

	dataMap := avalues.MergeMaps(
		avalues.PathValues,
		avalues.IconNames,
		avalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	dataMap[avalues.GrabResultStatusIconNameKey] = avalues.GrabResultStatusIconName(fileInfo.Status)
	dataMap[avalues.GrabResultItemHtmxOptionKey] = avalues.GrabResultStatusHtmxOption(fileInfo.Status, data.PathFileRow)
	dataMap[avalues.GrabResultItemStatusHtmlKey] = avalues.GrabResultStatusIconSvgRaw(fileInfo.Status, iconsDir)
	dataMap[avalues.GrabResultItemStatusTextKey] = fileInfo.StatusText
	if cacheChanged.youtubeTitle {
		dataMap[avalues.ResultYoutubeUrlFadeKey] = "fade-text"
	}
	if cacheChanged.FileSize {
		dataMap[avalues.ResultSizeFadeKey] = "fade-text"
	}
	if cacheChanged.Format {
		dataMap[avalues.ResultFormatFadeKey] = "fade-text"
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, dataMap)
	if err != nil {
		return nil, fasthttp.StatusInternalServerError, err
	}

	return &buf, fasthttp.StatusOK, nil
}
