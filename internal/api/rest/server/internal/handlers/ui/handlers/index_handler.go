package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

// IndexHandlers serves the main page (index.html)
func (h *DownloaderHandlers) IndexHandler(ctx *fasthttp.RequestCtx) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("authorization error", fasthttp.StatusUnauthorized, err))
		return
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	var rowsBuf bytes.Buffer
	err = h.getFilesHistory(ctx, &rowsBuf, userID, time.Now().UTC(), nil)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	systemInfo := h.usecases.Downloader.SystemInfo()

	showHistorySearch := h.usecases.Downloader.HistoryMode() != dtypes.HistoryModeDisabled

	cssPaths, err := uivalues.CssPaths(filepath.Join(h.assetsDir, dirStaticName, dirCssName))
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := uivalues.JsScripts(filepath.Join(h.assetsDir, dirStaticName, dirJsName))
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsImportMap, err := uivalues.JsImportMap(filepath.Join(h.assetsDir, dirStaticName, dirJsName))
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsImportMapJSON, err := json.MarshalIndent(
		map[string]any{
			"imports": jsImportMap,
		},
		"",
		"  ",
	)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	dataMap := uivalues.MergeMaps(uivalues.IndexValues, uivalues.FormGrabValues, uivalues.PathValues)
	dataMap[uivalues.CssPathsKey] = cssPaths
	dataMap[uivalues.JsScriptsKey] = jsScripts
	dataMap[uivalues.JsImportMapJSONKey] = template.HTML(jsImportMapJSON)
	dataMap[uivalues.ShowHistorySearchKey] = showHistorySearch
	dataMap[uivalues.ResultNoRowsKey] = rowsBuf.Len() == 0
	dataMap[uivalues.ResultRowsHTMLKey] = template.HTML(rowsBuf.String())
	dataMap[uivalues.AppVersionKey] = systemInfo.AppVersion
	dataMap[uivalues.DiskFreeKey] = uformat.BytesHuman(systemInfo.DiskFree)
	dataMap[uivalues.DiskUsedKey] = uformat.BytesHuman(systemInfo.DiskUsed)

	// Execute template with PageTitle
	if err := h.templates.ExecuteTemplate(ctx, uivalues.IndexHtmlFileName, dataMap); err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
		return
	}
}
