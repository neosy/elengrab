package handlers

import (
	"bytes"
	"html/template"
	"mime"
	"strconv"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/images"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	iconfig "github.com/neosy/elengrab/internal/config"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/humanize"
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

	cssPaths, err := paths.IndexCssPaths(h.assetFolders.Css())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := paths.IndexPageJsPaths(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	pwaManifestPath, err := paths.PwaManifestPath(h.assetFolders.Pwa())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var userAvatarActionMode = "none"
	if !h.downloader.DemoMode() {
		if ctxUser.UserType() < dtypes.UserTypeUser {
			userAvatarActionMode = "login"
		} else {
			userAvatarActionMode = "menu"
		}
	}

	imageData := &dtypes.ImageData{
		URL:    h.baseURL + paths.ImagePath(images.Elengrab1280ImageJpgFileName),
		Format: dtypes.ImageFormatJPEG,
		Width:  1280,
		Height: 720,
	}

	metaOgItems := make(pages.MetaOgItems, 0, 15)
	metaOgItems.Add("site_name", iconfig.AppName)
	metaOgItems.Add("type", "website")
	metaOgItems.Add("title", pages.PageTitle)
	metaOgItems.Add("description", pages.PageDescription)
	metaOgItems.Add("url", h.baseURL)
	metaOgItems.Add("image", imageData.URL)
	metaOgItems.Add("image:secure_url", imageData.URL)
	metaOgItems.Add("image:type", httpx.ContentTypeByExt(imageData.Format.String()))
	metaOgItems.Add("image:width", strconv.Itoa(imageData.Width))
	metaOgItems.Add("image:height", strconv.Itoa(imageData.Height))
	metaOgItems.Add("image:alt", "Elengrab logo")

	baseValues := pages.NewBaseValues()
	baseValues.MetaOgItems = metaOgItems

	extraData := make(map[string]any)
	extraData[items.UserAvatarIconKey] = icons.UserAvatarIconByType(ctxUser.UserType()).FileRaw()
	extraData[items.UserAvatarActionModeKey] = userAvatarActionMode
	extraData[items.ResultNoRowsKey] = rowsBuf.Len() == 0
	extraData[items.ResultRowsHTMLKey] = template.HTML(rowsBuf.String())

	pageData := pages.IndexPageData{
		BasePaths:  paths.NewPaths(),
		BaseValues: baseValues,
		Paths: pages.PagePaths{
			Css:         cssPaths,
			JsScripts:   jsScripts,
			PwaManifest: pwaManifestPath,
		},
		Values: pages.IndexPageValues{
			UserMenuSearchButtonIcon: icons.UserMenuSearchIcon.FileRaw(),
			SearchBackArrowIcon:      icons.SearchBackArrowIcon.FileRaw(),
			ShowHistorySearch:        true,
			DiskFree:                 humanize.Bytes(int64(systemInfo.DiskFree)),
			DiskUsed:                 humanize.Bytes(int64(systemInfo.DiskUsed)),
			GrabForm: pages.IndexGrabForm{
				InputPlaceholder:   pages.IndexGrabFormInputPlaceholder,
				SettingsButtonIcon: icons.IndexGrabSettingsButtonIcon.FileRaw(),
				GetButtonTitle:     pages.IndexGrabGetButtonTitle,
				GetButtonIcon:      icons.IndexGrabGetButtonIcon.FileRaw(),
			},
		},
		Extra: extraData,
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))

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
