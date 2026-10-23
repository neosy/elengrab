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
	requireSourceType                   dtypes.MediaDownloadSourceType
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
			URLTemplate: httppaths.DownloaderGroup + httppaths.MediaItemPath,
			Placeholder: RowMenuActionItemIDKey,
			NewTab:      false,
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
			URLTemplate: RowMenuActionURLKey,
			Placeholder: RowMenuActionURLKey,
			NewTab:      true,
		},

		icon:                icons.DownloaderRowMenuExternalLinkIcon,
		requireSourceType:   dtypes.MediaDownloadSourceTypeURL,
		requireEditAccess:   false,
		requireDeleteAccess: false,
	},

	{
		RenderType: renderTypeAction,
		Action:     rowMenuActionErrorInfo,
		Title:      "Error Information",

		icon:                icons.DownloaderRowMenuUpdateErrorInfoIcon,
		visibleStatuses:     []dtypes.MediaDownloadStatus{dtypes.MediaDownloadStatusFailed},
		requireDeleteAccess: true,
	},

	{
		RenderType: renderTypeDivider,
	},

	{
		RenderType: renderTypeAction,
		Action:     "share-link",
		Title:      "Share link",
		Link: linkOptions{
			URLTemplate: httppaths.DownloaderGroup + httppaths.MediaItemShortLinkPath,
			Placeholder: RowMenuActionItemIDKey,
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
			URLTemplate: httppaths.DownloaderGroup + httppaths.MediaItemShortLinkPath,
			Placeholder: RowMenuActionItemIDKey,
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
			URLTemplate: httppaths.DownloaderGroup + httppaths.MediaItemShortLinkPath,
			Placeholder: RowMenuActionItemIDKey,
		},

		icon:            icons.DownloaderRowMenuCopyLinkIcon,
		visibleStatuses: dtypes.MediaDownloadCompletedStatuses(),
	},

	{
		RenderType: renderTypeAction,
		Action:     rowMenuActionDeleteLink,
		Title:      "Delete short link",
		Link: linkOptions{
			URLTemplate: httppaths.DownloaderGroup + httppaths.MediaItemShortLinkPath,
			Placeholder: RowMenuActionItemIDKey,
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
			URLTemplate: httppaths.DownloaderGroup + httppaths.MediaItemRefreshPath,
			Placeholder: RowMenuActionItemIDKey,
			NewTab:      false,
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
			URLTemplate: httppaths.DownloaderGroup + httppaths.MediaItemEditPath,
			Placeholder: RowMenuActionItemIDKey,
			NewTab:      false,
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
			URLTemplate: httppaths.DownloaderGroup + httppaths.MediaItemPath,
			Placeholder: RowMenuActionItemIDKey,
		},

		icon:                icons.DownloaderRowMenuDeleteIcon,
		requireEditAccess:   false,
		requireDeleteAccess: true,
	},
}

func RowMenuActions(
	mapReplaceUrl map[string]string,
	status dtypes.MediaDownloadStatus,
	sourceType dtypes.MediaDownloadSourceType,
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

		if a.requireSourceType != dtypes.MediaDownloadSourceTypeNone && a.requireSourceType != sourceType {
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

		if key := action.Link.Placeholder; key != "" {
			value, ok := mapReplaceUrl[key]
			if ok {
				action.Link.URLTemplate = strings.Replace(action.Link.URLTemplate, key, value, 1)
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
