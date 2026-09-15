package menu

import (
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	RowMenuActionItemIDKey = "{itemId}"
	RowMenuActionURLKey    = "{url}"
)

type rowMenuAction struct {
	menuAction
	icon               icons.Icon
	visibleStatuses    []dtypes.MediaDownloadStatus
	disallowRefreshing bool

	requireVisibilityPublic             bool
	requireVisibilityPublicOrEditAccess bool
	requireEditAccess                   bool
	requireDeleteAccess                 bool
}

const (
	rowMenuActionErrorInfo  = "error-info"
	rowMenuActionCopyLink   = "copy-link"
	rowMenuActionCreateLink = "create-link"
	rowMenuActionDeleteLink = "delete-link"
	rowMenuActionEdit       = "edit"
)

var rowMenuActions = []rowMenuAction{
	{
		RenderType: renderTypeLink,
		Action:     "watch",
		Title:      "Watch",
		Link: linkOptions{
			URL:          httppaths.DownloaderGroup + httppaths.MediaItemWatchPath,
			NewTab:       false,
			replaceInURL: RowMenuActionItemIDKey,
		},

		icon:                icons.DownloaderRowMenuPlayIcon,
		visibleStatuses:     dtypes.MediaDownloadCompletedStatuses(),
		requireEditAccess:   false,
		requireDeleteAccess: false,
	},

	{
		RenderType: renderTypeLink,
		Action:     "open-original",
		Title:      "Open original in new tab",
		Link: linkOptions{
			URL:          RowMenuActionURLKey,
			NewTab:       true,
			replaceInURL: RowMenuActionURLKey,
		},

		icon:                icons.DownloaderRowMenuExternalLinkIcon,
		requireEditAccess:   false,
		requireDeleteAccess: false,
	},

	{
		RenderType: renderTypeAction,
		Action:     rowMenuActionErrorInfo,
		Title:      "Error Information",

		icon:                icons.DownloaderRowMenuUpdateErrorInfoIcon,
		visibleStatuses:     []dtypes.MediaDownloadStatus{dtypes.MediaDownloadStatusFailed},
		requireEditAccess:   true,
		requireDeleteAccess: false,
	},

	{
		RenderType: renderTypeDivider,
	},

	{
		RenderType: renderTypeAction,
		Action:     "share-link",
		Title:      "Share link",
		Link: linkOptions{
			URL:          httppaths.DownloaderGroup + httppaths.MediaItemShortLinkPath,
			replaceInURL: RowMenuActionItemIDKey,
		},

		icon:                                icons.DownloaderRowMenuShareLinkIcon,
		visibleStatuses:                     dtypes.MediaDownloadCompletedStatuses(),
		requireVisibilityPublicOrEditAccess: true,
	},

	{
		RenderType: renderTypeAction,
		Action:     rowMenuActionCreateLink,
		Title:      "Create short link",
		Link: linkOptions{
			URL:          httppaths.DownloaderGroup + httppaths.MediaItemShortLinkPath,
			replaceInURL: RowMenuActionItemIDKey,
		},

		icon:                                icons.DownloaderRowMenuCopyLinkIcon,
		visibleStatuses:                     dtypes.MediaDownloadCompletedStatuses(),
		requireVisibilityPublicOrEditAccess: true,
	},

	{
		RenderType: renderTypeAction,
		Action:     rowMenuActionCopyLink,
		Title:      "Copy short link",
		Link: linkOptions{
			URL:          httppaths.DownloaderGroup + httppaths.MediaItemShortLinkPath,
			replaceInURL: RowMenuActionItemIDKey,
		},

		icon:            icons.DownloaderRowMenuCopyLinkIcon,
		visibleStatuses: dtypes.MediaDownloadCompletedStatuses(),
	},

	{
		RenderType: renderTypeAction,
		Action:     rowMenuActionDeleteLink,
		Title:      "Delete short link",
		Link: linkOptions{
			URL:          httppaths.DownloaderGroup + httppaths.MediaItemShortLinkPath,
			replaceInURL: RowMenuActionItemIDKey,
		},

		icon:              icons.DownloaderRowMenuDeleteLinkIcon,
		visibleStatuses:   dtypes.MediaDownloadCompletedStatuses(),
		requireEditAccess: true,
	},

	{
		RenderType: renderTypeDivider,
	},

	{
		RenderType: renderTypeAction,
		Action:     "refresh",
		Title:      "Refresh Media Information",
		Link: linkOptions{
			URL:          httppaths.DownloaderGroup + httppaths.MediaItemRefreshPath,
			NewTab:       false,
			replaceInURL: RowMenuActionItemIDKey,
		},

		icon:                icons.DownloaderRowMenuUpdateMetadataIcon,
		visibleStatuses:     dtypes.MediaDownloadCompletedStatuses(),
		disallowRefreshing:  true,
		requireEditAccess:   true,
		requireDeleteAccess: false,
	},

	{
		RenderType: renderTypeLink,
		Action:     "edit",
		Title:      "Edit",
		Link: linkOptions{
			URL:          httppaths.DownloaderGroup + httppaths.MediaItemEditPath,
			NewTab:       false,
			replaceInURL: RowMenuActionItemIDKey,
		},

		icon:                icons.DownloaderRowMenuEditIcon,
		visibleStatuses:     dtypes.MediaDownloadEditableStatuses(),
		requireEditAccess:   true,
		requireDeleteAccess: false,
	},

	{
		RenderType: renderTypeDivider,
	},

	{
		RenderType: renderTypeAction,
		Action:     "delete",
		Title:      "Delete",
		Link: linkOptions{
			URL:          httppaths.DownloaderGroup + httppaths.MediaItemPath,
			replaceInURL: RowMenuActionItemIDKey,
		},

		icon:                icons.DownloaderRowMenuDeleteIcon,
		requireEditAccess:   false,
		requireDeleteAccess: true,
	},
}

func RowMenuActions(
	mapReplaceUrl map[string]string,
	status dtypes.MediaDownloadStatus,
	hasEditAccess bool,
	hasDeleteAccess bool,
	opts ...MenuActionOption,
) []rowMenuAction {
	actions := make([]rowMenuAction, 0, len(rowMenuActions))
	var (
		lastRenderType renderType
		options        = NewMenuActionOptions(opts...)
	)

	for _, a := range rowMenuActions {
		if len(a.visibleStatuses) > 0 && !status.Contains(a.visibleStatuses) {
			continue
		}

		if status == dtypes.MediaDownloadStatusRefreshing && a.disallowRefreshing {
			continue
		}

		if a.requireEditAccess && !hasEditAccess {
			continue
		}

		if a.requireDeleteAccess && !hasDeleteAccess {
			continue
		}

		if a.requireVisibilityPublic && options.Visibility != dtypes.MediaVisibilityPublic {
			continue
		}

		if a.requireVisibilityPublicOrEditAccess {
			if !hasEditAccess && options.Visibility != dtypes.MediaVisibilityPublic {
				continue
			}
		}

		switch a.Action {
		case rowMenuActionCreateLink:
			if options.HasShareLink {
				continue
			}
		case rowMenuActionCopyLink:
			if !options.HasShareLink {
				continue
			}
		case rowMenuActionDeleteLink:
			if !options.HasShareLink {
				continue
			}
		case rowMenuActionEdit:
			if !options.HasMetadata {
				continue
			}
		}

		if a.RenderType == renderTypeDivider {
			if len(actions) == 0 {
				continue
			}
			if lastRenderType == renderTypeDivider {
				continue
			}
		}

		lastRenderType = a.RenderType

		actions = append(actions, a)
		action := &actions[len(actions)-1]

		if action.icon.FileName() != "" {
			action.IconSvg = action.icon.FileRaw()
		}

		if key := action.Link.replaceInURL; key != "" {
			value, ok := mapReplaceUrl[key]
			if ok {
				action.Link.URL = strings.Replace(action.Link.URL, key, value, 1)
			}
		}

		if action.Action == rowMenuActionErrorInfo && options.ErrorText != "" {
			action.Text = options.ErrorText
		}
	}

	if len(actions) > 0 && actions[len(actions)-1].RenderType == renderTypeDivider {
		actions = actions[:len(actions)-1]
	}

	return actions
}
