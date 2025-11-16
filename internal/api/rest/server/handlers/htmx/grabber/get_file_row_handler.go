package grabberh

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/google/uuid"
	avalues "github.com/neosy/elengrab/internal/api/rest/server/assets/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
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
	DownloadURL  string
}

type cacheRowEntry struct {
	youtubeTitle string
	FileSize     int
	Format       string
	Updated      time.Time
}

type cache[T any] struct {
	data map[string]T
	mu   sync.RWMutex
	ttl  time.Duration
}

var (
	cacheRow = cache[cacheRowEntry]{
		data: make(map[string]cacheRowEntry),
		ttl:  time.Minute,
	}
)

func (h *GrabberHandlers) GetFileRow(ctx *fasthttp.RequestCtx) {
	var cacheChanged = struct {
		youtubeTitle bool
		FileSize     bool
		Format       bool
	}{}

	var (
		tmplPath string
		data     fileRowInfoData
	)

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is required")
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is incorrect")
		return
	}

	resp, err := h.usecases.Downloader.GetFileInfo(ctx, fileId)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(fmt.Sprintf("Internal error: %v", err))
		return
	}

	if resp == nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("the request returned an empty")
		return
	}

	// Checking the cache
	{
		path := string(ctx.Path())

		cacheRow.mu.RLock()
		cached, exists := cacheRow.data[path]
		cacheRow.mu.RUnlock()

		if exists {
			cacheChanged.youtubeTitle = cached.youtubeTitle != resp.YoutubeTitle
			if resp.FileSize != nil {
				cacheChanged.FileSize = cached.FileSize != *resp.FileSize
			}
			cacheChanged.Format = cached.Format != resp.FileExt
		}

		if exists && resp.Status != dtypes.FileStatusDone &&
			!cacheChanged.youtubeTitle &&
			!cacheChanged.FileSize &&
			!cacheChanged.Format &&
			time.Since(cached.Updated) < cacheRow.ttl {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			ctx.Response.Header.Set("HX-Trigger", "no-op")
			return
		}

		// Updating the cache
		var fileSize int
		if resp.FileSize != nil {
			fileSize = *resp.FileSize
		}
		cacheRow.mu.Lock()
		cacheRow.data[path] = cacheRowEntry{
			youtubeTitle: resp.YoutubeTitle,
			FileSize:     fileSize,
			Format:       resp.FileExt,
			Updated:      time.Now(),
		}
		cacheRow.mu.Unlock()
	}

	tmplPath = filepath.Join(h.assetsDir, "templates", avalues.GrabResultStatusHtmlFileName(resp.Status))

	data.YoutubeURL = resp.YoutubeUrl
	data.YoutubeTitle = resp.YoutubeTitle
	data.PathFileRow = httppaths.BuildPathFileRow(resp.FileId)

	data.FileSize = "-"
	if resp.FileSize != nil && *resp.FileSize > 0 {
		data.FileSize = utils.BytesToHuman(*resp.FileSize)
	}

	data.Format = "-"
	if resp.FileExt != "" {
		data.Format = resp.FileExt
	}
	// Set URL for download endpoint
	data.DownloadURL = httppaths.BuildPathFileDownload(resp.FileId)

	tpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(fmt.Sprintf("Template error: %v", err))
		return
	}

	dataMap := avalues.MergeMaps(
		avalues.PathValues,
		avalues.IconNames,
		avalues.StructToMap(data),
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	dataMap[avalues.GrabResultStatusIconNameKey] = avalues.GrabResultStatusIconName(resp.Status)
	dataMap[avalues.GrabResultItemHtmxOptionKey] = avalues.GrabResultStatusHtmxOption(resp.Status, data.PathFileRow)
	dataMap[avalues.GrabResultItemStatusHtmlKey] = avalues.GrabResultStatusIconSvgRaw(resp.Status, iconsDir)
	dataMap[avalues.GrabResultItemStatusTextKey] = resp.StatusText
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
	tpl.Execute(&buf, dataMap)
	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
