package downloader

import (
	"bytes"
	"html/template"
	"mime"
	"strconv"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/clientcap"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/images"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	iconfig "github.com/neosy/elengrab/internal/config"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/humanize"
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

	systemInfo := h.downloader.SystemInfo()

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

	var (
		userAvatarActionMode = "none"
		userMenuAvatarTitle  = ""
	)
	if !h.downloader.DemoMode() {
		if authCtx.UserType() < dtypes.UserTypeUser {
			userAvatarActionMode = "login"
			userMenuAvatarTitle = "Login"
		} else {
			userAvatarActionMode = "menu"
			userMenuAvatarTitle = "Account menu"
		}
	}

	imageData := &dtypes.ImageData{
		URL:    h.baseURL + paths.ImagePath(images.Elengrab1280ImageJpgFileName),
		Format: dtypes.ImageFormatJPEG,
		Width:  1280,
		Height: 720,
	}

	queryFilters := h.mappers.MapQueryFiltersDomainToFilters(query.Filters)

	searchParameters := types.NewSearchParameters()
	searchParameters.AddValues(query.ViewMode, queryFilters, dtypes.QueryMediaDownloadCursor{})

	channelHeaderPageData := pages.ChannelHeader{}
	if channelID := searchParameters.FindChannelID(); channelID != uuid.Nil {
		channelHeaderPageData = h.buildChannelHeaderPageData(ctx, channelID)
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
	extraData[items.UserAvatarIconKey] = icons.UserAvatarIconByType(authCtx.UserType()).FileRaw()
	extraData[items.UserAvatarActionModeKey] = userAvatarActionMode

	pageData := pages.ChannelPageData{
		BasePaths:  paths.NewHttpPaths(),
		BaseValues: baseValues,
		Paths: pages.PagePaths{
			Css:         cssPaths,
			JsScripts:   jsScripts,
			PwaManifest: pwaManifestPath,
		},
		Values: pages.ChannelPageValues{
			HeaderActionsSearchButtonIcon:   icons.UserMenuSearchIcon.FileRaw(),
			SearchBackArrowIcon:             icons.SearchBackArrowIcon.FileRaw(),
			HeaderActionsDownloadButtonIcon: icons.UserMenuDownloadIcon.FileRaw(),
			ShowHistorySearch:               true,
			HeaderActionsAvatarTitle:        userMenuAvatarTitle,
			HasCreateAccess:                 h.downloader.CanCreateMediaDownload(authCtx),
			HasWriteOperationAccess:         h.downloader.HasWriteOperationAccess(authCtx),
			DiskFree:                        humanize.Bytes(int64(systemInfo.DiskFree)),
			DiskUsed:                        humanize.Bytes(int64(systemInfo.DiskUsed)),

			SearchQuery: query.GetSearchQueryString(),

			ChannelHeader: channelHeaderPageData,
			ChannelJSON:   string(channelHeaderPageData.JSON()),

			ActiveViewMode: query.ViewMode.String(),
			ViewModeTabs:   pages.BuildViewModeTabs(query.ViewMode),

			HasSearchFilters:     searchParameters.QueryFilters().Len() != 0,
			SearchParametersJSON: string(searchParameters.BuildJSON()),
			SearchParameters:     template.HTML(searchParameters.EncodeShortQueryValue()),

			ResultNoRows:   rowsBuf.Len() == 0,
			ResultRowsHTML: template.HTML(rowsBuf.String()),

			VideoPreview: pages.VideoPreview{
				SoundOnIcon:  icons.VideoPreviewSoundOnIcon.FileRaw(),
				SoundOffIcon: icons.VideoPreviewSoundOffIcon.FileRaw(),
			},
		},
		Extra: extraData,
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))

	// Execute template with PageTitle
	if err := h.templates.Pages[pages.ChannelPage.Key()].ExecuteTemplate(ctx, pages.ChannelPage.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
