package handlers

import (
	"bytes"
	"html/template"
	"path/filepath"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/images"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	iconfig "github.com/neosy/elengrab/internal/config"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) renderIndexPage(ctx *fasthttp.RequestCtx, ctxUser *dauth.UserContext) {

	var rowsBuf bytes.Buffer
	err := h.getDownloadsHistory(ctx, &rowsBuf, *ctxUser, time.Now().UTC(), nil)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	systemInfo := h.downloader.SystemInfo()

	cssPaths, err := paths.CssIndexPaths(h.assetFolders.Css())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := paths.JsIndexPaths(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsImportMapJSON, err := paths.JsIndexImportMapJSON(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	pwaManifestPath, err := paths.PwaManifestPath(h.assetFolders.Pwa())
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

	imageURL := h.baseURL + paths.ImagePath(images.Elengrab1280ImageJpgFileName)

	metaOgItems := make(pages.MetaOgItems, 0, 15)
	metaOgItems.Add("site_name", iconfig.AppName)
	metaOgItems.Add("type", "website")
	metaOgItems.Add("title", pages.PageTitle)
	metaOgItems.Add("description", pages.PageDescription)
	metaOgItems.Add("url", h.baseURL)
	metaOgItems.Add("image", imageURL)
	metaOgItems.Add("image:secure_url", imageURL)
	metaOgItems.Add("image:type", httpx.ContentTypeFromPath(imageURL))
	metaOgItems.Add("image:width", "1280")
	metaOgItems.Add("image:height", "720")
	metaOgItems.Add("image:alt", "Elengrab logo")

	baseValues := pages.NewBaseValues()
	baseValues.MetaOgItems = metaOgItems

	extraData := make(map[string]any)
	extraData[items.UserAvatarIconKey] = icons.FileRawByKey(icons.UserAvatarKeyByType(ctxUser.UserType()), iconsDir)
	extraData[items.UserAvatarActionModeKey] = userAvatarActionMode
	extraData[items.ResultNoRowsKey] = rowsBuf.Len() == 0
	extraData[items.ResultRowsHTMLKey] = template.HTML(rowsBuf.String())

	pageData := pages.IndexPageData{
		BasePaths:  paths.NewPaths(),
		BaseValues: baseValues,
		Paths: pages.PagePaths{
			Css:             cssPaths,
			JsScripts:       jsScripts,
			PwaManifest:     pwaManifestPath,
			JsImportMapJSON: template.HTML(jsImportMapJSON),
		},
		Values: pages.IndexPageValues{
			ShowHistorySearch: true,
			DiskFree:          uformat.BytesHuman(int64(systemInfo.DiskFree)),
			DiskUsed:          uformat.BytesHuman(int64(systemInfo.DiskUsed)),
			GrabForm: pages.IndexGrabForm{
				InputPlaceholder: pages.IndexGrabFormInputPlaceholder,
			},
		},
		Extra: extraData,
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	// Load template
	tmpl, err := h.loadPageTemplate(pages.IndexPage.FileName())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template with PageTitle
	if err := tmpl.ExecuteTemplate(ctx, pages.IndexPage.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
