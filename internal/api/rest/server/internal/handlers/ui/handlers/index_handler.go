package handlers

import (
	"bytes"
	"html/template"
	"path/filepath"
	"time"

	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	iconfig "github.com/neosy/elengrab/internal/config"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

// IndexHandlers serves the main page (index.html)
func (h *DownloaderHandlers) IndexHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsHead() {
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	if h.redirectGuestIfAuthRequired(ctx) {
		return
	}

	ctxUser := authmw.UserFromContext(ctx)
	if ctxUser == nil {
		// anonymous
		ctxUser = dauth.UserContextAnonymous()
	}

	var rowsBuf bytes.Buffer
	err := h.getFilesHistory(ctx, &rowsBuf, *ctxUser, time.Now().UTC(), nil)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	systemInfo := h.downloader.SystemInfo()

	cssPaths, err := uivalues.CssIndexPaths(h.assetFolders.Css())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := uivalues.JsIndexPaths(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsImportMapJSON, err := uivalues.JsIndexImportMapJSON(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	var userAvatarActionMode = "none"
	if !h.downloader.DemoMode() {
		if ctxUser.UserType() < dtypes.UserTypeUser {
			userAvatarActionMode = "login"
		} else {
			userAvatarActionMode = "menu"
		}
	}

	imageURL := h.baseURL + uivalues.ImageHttpPath(uivalues.Elengrab1280ImageFileName)

	metaOgItems := make(uivalues.MetaOgItems, 0, 15)
	metaOgItems.Add("site_name", iconfig.AppName)
	metaOgItems.Add("title", uivalues.PageTitle)
	metaOgItems.Add("description", uivalues.PageDescription)
	metaOgItems.Add("url", h.baseURL)
	metaOgItems.Add("image", imageURL)
	metaOgItems.Add("image:width", "1280")
	metaOgItems.Add("image:height", "720")

	baseValues := uivalues.NewBaseValues()
	baseValues.MetaOgItems = metaOgItems

	dataMap := uivalues.MergeMaps(
		baseValues.ToMap(),
		uivalues.FormGrabValues,
		uivalues.PathValues,
	)
	dataMap[uivalues.CssPathsKey] = cssPaths
	dataMap[uivalues.JsScriptsKey] = jsScripts
	dataMap[uivalues.JsImportMapJSONKey] = template.HTML(jsImportMapJSON)
	dataMap[uivalues.UserAvatarIconKey] = template.HTML(
		uivalues.IconFileRawByKey(uivalues.UserAvatarKeyByType(ctxUser.UserType()), iconsDir))
	dataMap[uivalues.UserAvatarActionModeKey] = userAvatarActionMode
	dataMap[uivalues.ShowHistorySearchKey] = true
	dataMap[uivalues.ResultNoRowsKey] = rowsBuf.Len() == 0
	dataMap[uivalues.ResultRowsHTMLKey] = template.HTML(rowsBuf.String())
	dataMap[uivalues.AppVersionKey] = systemInfo.AppVersion
	dataMap[uivalues.DiskFreeKey] = uformat.BytesHuman(systemInfo.DiskFree)
	dataMap[uivalues.DiskUsedKey] = uformat.BytesHuman(systemInfo.DiskUsed)

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	// Load template
	tmpl, err := h.loadPage(uivalues.PageIndex.FileName())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template with PageTitle
	if err := tmpl.ExecuteTemplate(ctx, uivalues.PageIndex.Key(), dataMap); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
