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
	requireWriteAccess bool
}

const (
	rowMenuActionErrorInfo  = "error-info"
	rowMenuActionCopyLink   = "copy-link"
	rowMenuActionCreateLink = "create-link"
	rowMenuActionDeleteLink = "delete-link"
)

var rowMenuActions = []rowMenuAction{
	{
		menuAction: menuAction{
			RenderType: renderTypeLink,
			Action:     "watch",
			Title:      "Watch",
			Link: linkOptions{
				URL:          httppaths.GroupDownloader + httppaths.PathMediaItemWatch,
				NewTab:       false,
				replaceInURL: RowMenuActionItemIDKey,
			},
		},
		icon:               icons.DownloaderRowMenuPlayIcon,
		visibleStatuses:    dtypes.MediaDownloadCompletedStatuses(),
		requireWriteAccess: false,
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeLink,
			Action:     "open-original",
			Title:      "Open original in new tab",
			Link: linkOptions{
				URL:          RowMenuActionURLKey,
				NewTab:       true,
				replaceInURL: RowMenuActionURLKey,
			},
		},
		icon:               icons.DownloaderRowMenuExternalLinkIcon,
		requireWriteAccess: false,
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     rowMenuActionErrorInfo,
			Title:      "Error Information",
		},
		icon:               icons.DownloaderRowMenuUpdateErrorInfoIcon,
		visibleStatuses:    []dtypes.MediaDownloadStatus{dtypes.MediaDownloadStatusFailed},
		requireWriteAccess: true,
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeDivider,
		},
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     "share-link",
			Title:      "Share link",
			Link: linkOptions{
				URL:          httppaths.GroupDownloader + httppaths.PathMediaItemShortLink,
				replaceInURL: RowMenuActionItemIDKey,
			},
		},
		icon:               icons.DownloaderRowMenuShareLinkIcon,
		visibleStatuses:    dtypes.MediaDownloadCompletedStatuses(),
		requireWriteAccess: true,
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     rowMenuActionCreateLink,
			Title:      "Create short link",
			Link: linkOptions{
				URL:          httppaths.GroupDownloader + httppaths.PathMediaItemShortLink,
				replaceInURL: RowMenuActionItemIDKey,
			},
		},
		icon:               icons.DownloaderRowMenuCopyLinkIcon,
		visibleStatuses:    dtypes.MediaDownloadCompletedStatuses(),
		requireWriteAccess: true,
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     rowMenuActionCopyLink,
			Title:      "Copy short link",
			Link: linkOptions{
				URL:          httppaths.GroupDownloader + httppaths.PathMediaItemShortLink,
				replaceInURL: RowMenuActionItemIDKey,
			},
		},
		icon:               icons.DownloaderRowMenuCopyLinkIcon,
		visibleStatuses:    dtypes.MediaDownloadCompletedStatuses(),
		requireWriteAccess: false,
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     rowMenuActionDeleteLink,
			Title:      "Delete short link",
			Link: linkOptions{
				URL:          httppaths.GroupDownloader + httppaths.PathMediaItemShortLink,
				replaceInURL: RowMenuActionItemIDKey,
			},
		},
		icon:               icons.DownloaderRowMenuDeleteLinkIcon,
		visibleStatuses:    dtypes.MediaDownloadCompletedStatuses(),
		requireWriteAccess: true,
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeDivider,
		},
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     "refresh",
			Title:      "Refresh Media Information",
			Link: linkOptions{
				URL:          httppaths.GroupDownloader + httppaths.PathMediaItemRefresh,
				NewTab:       false,
				replaceInURL: RowMenuActionItemIDKey,
			},
		},
		icon:               icons.DownloaderRowMenuUpdateMetadataIcon,
		visibleStatuses:    dtypes.MediaDownloadCompletedStatuses(),
		disallowRefreshing: true,
		requireWriteAccess: true,
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeLink,
			Action:     "edit",
			Title:      "Edit",
			Link: linkOptions{
				URL:          httppaths.GroupDownloader + httppaths.PathMediaItemEdit,
				NewTab:       false,
				replaceInURL: RowMenuActionItemIDKey,
			},
		},
		icon:               icons.DownloaderRowMenuEditIcon,
		visibleStatuses:    dtypes.MediaDownloadCompletedStatuses(),
		requireWriteAccess: true,
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeDivider,
		},
	},

	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     "delete",
			Title:      "Delete",
			Link: linkOptions{
				URL:          httppaths.GroupDownloader + httppaths.PathMediaItem,
				replaceInURL: RowMenuActionItemIDKey,
			},
		},
		icon:               icons.DownloaderRowMenuDeleteIcon,
		requireWriteAccess: true,
	},
}

func RowMenuActions(
	mapReplaceUrl map[string]string,
	status dtypes.MediaDownloadStatus,
	hasWriteAccess bool,
	opts ...MenuActionOption,
) []rowMenuAction {
	actions := make([]rowMenuAction, 0, len(rowMenuActions))
	var (
		lastRenderType renderType
		options        MenuActionOptions
	)

	for _, opt := range opts {
		opt(&options)
	}

	for _, a := range rowMenuActions {
		if len(a.visibleStatuses) > 0 && !status.Contains(a.visibleStatuses) {
			continue
		}

		if status == dtypes.MediaDownloadStatusRefreshing && a.disallowRefreshing {
			continue
		}

		if a.requireWriteAccess && !hasWriteAccess {
			continue
		}

		if a.Action == rowMenuActionCreateLink && options.HasShareLink {
			continue
		}

		if a.Action == rowMenuActionCopyLink && !options.HasShareLink {
			continue
		}

		if a.Action == rowMenuActionDeleteLink && !options.HasShareLink {
			continue
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
