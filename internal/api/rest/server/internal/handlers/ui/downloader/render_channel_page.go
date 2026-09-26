package downloader

import (
	"bytes"
	"html/template"
	"mime"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/clientcap"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	pagesdata "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/pages_data"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) renderChannelPage(
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

	cssPaths, err := h.assetPaths.ChannelPageCssPaths()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	caps := clientcap.Detect(string(ctx.UserAgent()))

	jsScripts, err := h.assetPaths.ChannelPageJsPaths(caps.IsLegacyWebKit)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	pwaManifestPath, err := h.assetPaths.PwaManifestPath()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	queryFilters := h.mappers.MapQueryFiltersDomainToFilters(query.Filters)

	channelHeaderPageData := pages.ChannelHeader{}
	if filterChannelID, exists := queryFilters.Find(qkeys.ChannelIDKey); exists {
		channelID, _ := idcodec.DecodeUUIDBase64URL(filterChannelID.Value)
		if channelID != uuid.Nil {
			channelHeaderPageData = h.buildChannelHeaderPageData(ctx, channelID)
		}
	}

	baseValues := pages.NewBaseValues()
	baseValues.MetaOgItems = pagesdata.BuildMetaOgItems(h.baseURL)

	pagesListValues := pagesdata.BuildPagesListValues(authCtx, query, queryFilters, h.downloader)
	pagesListValues.ResultNoRows = rowsBuf.Len() == 0
	pagesListValues.ResultRowsHTML = template.HTML(rowsBuf.String())

	pageData := pages.ChannelPageData{
		BasePaths:  paths.NewHttpPaths(),
		BaseValues: baseValues,
		Paths: pages.PagePaths{
			Css:         cssPaths,
			JsScripts:   jsScripts,
			PwaManifest: pwaManifestPath,
		},
		Values: pages.ChannelPageValues{
			PagesListValues: pagesListValues,

			ChannelHeader: channelHeaderPageData,
		},
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))

	// Execute template with PageTitle
	if err := h.templates.Pages[pages.ChannelPage.Key()].ExecuteTemplate(ctx, pages.ChannelPage.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
