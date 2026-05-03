package handlers

import (
	"bytes"
	"html/template"
	"path/filepath"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
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

	ctxUser := policy.ResolveUserOrAnonym(ctx)

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

	imageURL := h.baseURL + uivalues.ImageHttpPath(uivalues.Elengrab1280ImageJpgFileName)

	metaOgItems := make(uivalues.MetaOgItems, 0, 15)
	metaOgItems.Add("site_name", iconfig.AppName)
	metaOgItems.Add("type", "website")
	metaOgItems.Add("title", uivalues.PageTitle)
	metaOgItems.Add("description", uivalues.PageDescription)
	metaOgItems.Add("url", h.baseURL)
	metaOgItems.Add("image", imageURL)
	metaOgItems.Add("image:secure_url", imageURL)
	metaOgItems.Add("image:type", httpx.ContentTypeFromPath(imageURL))
	metaOgItems.Add("image:width", "1280")
	metaOgItems.Add("image:height", "720")
	metaOgItems.Add("image:alt", "Elengrab logo")

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
	dataMap[uivalues.DiskFreeKey] = uformat.BytesHuman(int64(systemInfo.DiskFree))
	dataMap[uivalues.DiskUsedKey] = uformat.BytesHuman(int64(systemInfo.DiskUsed))

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
