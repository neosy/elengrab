package pagesdata

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func BuildPagesListValues(
	authCtx dauth.AuthContext,
	query udto.MediaDownloadQuery,
	queryFilters *types.QueryFilters,
	downloader downloader.DownloaderAPI,
) pages.PagesListValues {
	var (
		userAvatarActionMode = "none"
		userMenuAvatarTitle  = ""
	)
	if !downloader.DemoMode() {
		if authCtx.UserType() < dtypes.UserTypeUser {
			userAvatarActionMode = "login"
			userMenuAvatarTitle = "Login"
		} else {
			userAvatarActionMode = "menu"
			userMenuAvatarTitle = "Account menu"
		}
	}

	searchParam := types.NewSearchParameters()
	searchParam.AddValues(query.ViewMode, queryFilters, dtypes.QueryMediaDownloadCursor{})

	searchFilters := queryFilters.FilterByKeys(qkeys.SearchFilterKeys)
	searchParamFilters := searchParam.QueryFilters()

	return pages.PagesListValues{
		HeaderActionsSearchButtonIcon:   icons.UserMenuSearchIcon.FileRaw(),
		HeaderActionsDownloadButtonIcon: icons.UserMenuDownloadIcon.FileRaw(),
		HeaderActionsAvatarTitle:        userMenuAvatarTitle,

		UserAvatar: pages.UserAvatar{
			Icon:       icons.UserAvatarIconByType(authCtx.UserType()).FileRaw(),
			ActionMode: userAvatarActionMode,
		},

		ShowHistorySearch:   true,
		SearchBackArrowIcon: icons.SearchBackArrowIcon.FileRaw(),

		HasCreateAccess:         downloader.CanCreateMediaDownload(authCtx),
		HasWriteOperationAccess: downloader.HasWriteOperationAccess(authCtx),

		ActiveViewMode: query.ViewMode.String(),
		ViewModeTabs:   pages.BuildViewModeTabs(query.ViewMode),

		SearchQuery: query.GetSearchQueryString(),

		HasSearchFilters:     searchFilters.Len() != 0 || searchParamFilters.Len() != 0,
		SearchFiltersJSON:    string(searchFilters.BuildJSON()),
		SearchParametersJSON: string(searchParam.BuildJSON()),
		SearchParameters:     template.HTML(searchParam.EncodeShortQueryValue()),

		VideoPreview: pages.VideoPreview{
			SoundOnIcon:  icons.VideoPreviewSoundOnIcon.FileRaw(),
			SoundOffIcon: icons.VideoPreviewSoundOffIcon.FileRaw(),
			PlayIcon:     icons.VideoPreviewPlayIcon.FileRaw(),
			PauseIcon:    icons.VideoPreviewPauseIcon.FileRaw(),
		},
	}
}
