package downloader

import (
	"bytes"
	"html/template"
	"mime"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/clientcap"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	pagesdata "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/pages_data"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) renderIndexPage(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	query udto.MediaDownloadQuery,
) {
	var rowsBuf bytes.Buffer
	err := h.renderDownloadItemList(ctx, &rowsBuf, authCtx, query)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	cssPaths, err := h.assetPaths.IndexPageCssPaths()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	caps := clientcap.Detect(string(ctx.UserAgent()))

	jsScripts, err := h.assetPaths.IndexPageJsPaths(caps.IsLegacyWebKit)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	pwaManifestPath, err := h.assetPaths.PwaManifestPath()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	systemInfo := h.downloader.SystemInfo()

	baseValues := pages.NewBaseValues()
	baseValues.MetaOgItems = pagesdata.BuildMetaOgItems(h.baseURL)

	queryFilters := h.mappers.MapQueryFiltersDomainToFilters(query.Filters)

	pagesListValues := pagesdata.BuildPagesListValues(authCtx, query, queryFilters, h.downloader)
	pagesListValues.ResultNoRows = rowsBuf.Len() == 0
	pagesListValues.ResultRowsHTML = template.HTML(rowsBuf.String())
	pagesListValues.AboutDialog = h.mappers.MapSystemInfoToAboutDialogValues(systemInfo)

	pageData := pages.IndexPageData{
		BasePaths:  paths.NewHttpPaths(),
		BaseValues: baseValues,
		Paths: pages.PagePaths{
			Css:         cssPaths,
			JsScripts:   jsScripts,
			PwaManifest: pwaManifestPath,
		},
		Values: pages.IndexPageValues{
			PagesListValues: pagesListValues,

			GrabForm: pages.IndexGrabForm{
				InputPlaceholder:   pages.IndexGrabFormInputPlaceholder,
				SettingsButtonIcon: icons.IndexGrabSettingsButtonIcon.FileRaw(),
				GetButtonTitle:     pages.IndexGrabGetButtonTitle,
				GetButtonIcon:      icons.IndexGrabGetButtonIcon.FileRaw(),
			},

			DiskFree: humanize.Bytes(int64(systemInfo.DiskFree)),
			DiskUsed: humanize.Bytes(int64(systemInfo.DiskUsed)),
		},
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))

	// Execute template with PageTitle
	if err := h.templates.Pages[pages.IndexPage.Key()].ExecuteTemplate(ctx, pages.IndexPage.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
