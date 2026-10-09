package downloader

import (
	"bytes"
	"html/template"
	"mime"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/clientcap"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	pagesdata "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/pages_data"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) renderUserDownloadsPage(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	user dto.AuthUserResponse,
	query udto.MediaDownloadQuery,
) {
	var rowsBuf bytes.Buffer
	err := h.renderDownloadItemList(ctx, &rowsBuf, authCtx, query)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	cssPaths, err := h.assetPaths.UserDownloadsPageCssPaths()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	caps := clientcap.Detect(string(ctx.UserAgent()))

	jsScripts, err := h.assetPaths.UserDownloadsPageJsPaths(caps.IsLegacyWebKit)
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

	queryFilters, err = h.prepareSearchQueryFilters(ctx, queryFilters)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var userLogin string

	if queryFilters != nil {
		if filter, exists := queryFilters.Find(qkeys.UserNameKey); exists {
			userLogin = filter.Value
		}
	}

	userHeaderPageData := h.buildUserDownloadsHeaderPageData(ctx, userLogin)

	userDisplayName := userDisplayName(
		&user.UserID,
		user.Login,
		user.UserType(),
		authCtx.UserID,
	)

	if userDisplayName != "" {
		userHeaderPageData.Title = userDisplayName
	}

	if query.Filters != nil {
		if filter, exists := query.Filters.Find(dtypes.QueryFilterNameUserID); exists {
			userID, ok := filter.Value().(uuid.UUID)
			if ok && userID == authCtx.UserID {
				userHeaderPageData.Title = "My Downloads"
			}
		}
	}

	systemInfo := h.downloader.SystemInfo()

	baseValues := pages.NewBaseValues()
	baseValues.MetaOgItems = pagesdata.BuildMetaOgItems(h.baseURL)

	pagesListValues := pagesdata.BuildPagesListValues(authCtx, query, queryFilters, h.downloader)
	pagesListValues.ResultNoRows = rowsBuf.Len() == 0
	pagesListValues.ResultRowsHTML = template.HTML(rowsBuf.String())
	pagesListValues.AboutDialog = h.mappers.MapSystemInfoToAboutDialogValues(systemInfo)
	pagesListValues.AudioPlayingIcon = icons.MediaAudioPlayingIcon.FileRaw()
	pagesListValues.AudioPlayIcon = icons.MediaAudioPlayIcon.FileRaw()
	pagesListValues.AudioPauseIcon = icons.MediaAudioPauseIcon.FileRaw()

	pageData := pages.UserDownloadsPageData{
		BasePaths:  paths.NewHttpPaths(),
		BaseValues: baseValues,
		Paths: pages.PagePaths{
			Css:         cssPaths,
			JsScripts:   jsScripts,
			PwaManifest: pwaManifestPath,
		},
		Values: pages.UserDownloadsPageValues{
			PagesListValues: pagesListValues,

			UserHeader: userHeaderPageData,
		},
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))

	// Execute template with PageTitle
	if err := h.templates.Pages[pages.UserDownloadsPage.Key()].ExecuteTemplate(ctx, pages.UserDownloadsPage.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
