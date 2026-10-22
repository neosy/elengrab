package downloader

import (
	"mime"

	"uuid"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/clientcap"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

type renderEditMediaPageRequest struct {
	pageURL    string
	downloadID uuid.UUID

	authCtx dauth.AuthContext
}

func (h *DownloaderHandlers) renderEditMediaPage(
	ctx *fasthttp.RequestCtx,
	req renderEditMediaPageRequest,
) {
	if ctx.IsHead() {
		fasthttpx.WriteHTMLHeadResponse(ctx)
		return
	}

	downloadInfo, err := h.downloader.GetDownloadInfoForEdit(ctx, req.authCtx, req.downloadID)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	baseValues := pages.NewBaseValues()
	baseValues.Title = downloadInfo.MediaTitle

	cssPaths, err := h.assetPaths.EditMediaPageCssPaths()
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	caps := clientcap.Detect(string(ctx.UserAgent()))

	jsScripts, err := h.assetPaths.EditMediaPageJsPaths(caps.IsLegacyWebKit)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	pwaManifestPath, err := h.assetPaths.PwaManifestPath()
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	pageData := pages.EditMediaPageData{
		BaseValues: baseValues,
		BasePaths:  paths.NewHttpPaths(),
		Paths: pages.PagePaths{
			Css:         cssPaths,
			JsScripts:   jsScripts,
			PwaManifest: pwaManifestPath,
		},
		Values: pages.EditMediaPageValues{
			MediaTitle:       downloadInfo.MediaTitle,
			MediaDescription: downloadInfo.MediaDescription,
			PatchURL:         httppaths.BuildMediaItemPath(downloadInfo.DownloadID),
			Visibility:       downloadInfo.Visibility.String(),
			VisibilityList:   pages.VisibilityList(),
		},
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))

	// Execute template with PageTitle
	if err := h.templates.Pages[pages.EditMediaPage.Key()].ExecuteTemplate(ctx, pages.EditMediaPage.Key(), pageData); err != nil {
		fasthttpx.WriteErrorx(ctx, errInternal(err))
		return
	}
}
