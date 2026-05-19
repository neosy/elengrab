package menu

import (
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

type rowMenuAction struct {
	RenderType string `json:"renderType"` // "link" | "action" | "divider"
	Action     string `json:"action"`
	Title      string `json:"title"`
	IconSvg    any    `json:"iconSvg,omitempty"`

	icon           icons.Icon
	URL            string
	NewTab         bool
	replaceInURL   string
	onlyStatusDone bool
}

const (
	menuActionRenderTypeLink    = "link"
	menuActionRenderTypeAction  = "action"
	menuActionRenderTypeDivider = "divider"
)

const (
	RowMenuActionItemIDKey = "{itemId}"
	RowMenuActionURLKey    = "{url}"
)

var rowMenuActions = []rowMenuAction{
	{
		RenderType:     menuActionRenderTypeLink,
		Action:         "watch",
		Title:          "Watch",
		icon:           icons.DownloaderRowMenuPlayIcon,
		URL:            httppaths.GroupDownloader + httppaths.PathMediaItemWatch,
		NewTab:         false,
		replaceInURL:   RowMenuActionItemIDKey,
		onlyStatusDone: true,
	},
	{
		RenderType:   menuActionRenderTypeLink,
		Action:       "open-original",
		Title:        "Open original in new tab",
		icon:         icons.DownloaderRowMenuExternalLinkIcon,
		URL:          RowMenuActionURLKey,
		NewTab:       true,
		replaceInURL: RowMenuActionURLKey,
	},

	{RenderType: menuActionRenderTypeDivider},

	{
		RenderType:     menuActionRenderTypeAction,
		Action:         "share-link",
		Title:          "Share link",
		icon:           icons.DownloaderRowMenuShareLinkIcon,
		URL:            httppaths.GroupDownloader + httppaths.PathMediaItemShortLink,
		replaceInURL:   RowMenuActionItemIDKey,
		onlyStatusDone: true,
	},

	{
		RenderType:     menuActionRenderTypeAction,
		Action:         "copy-link",
		Title:          "Copy short link",
		icon:           icons.DownloaderRowMenuCopyLinkIcon,
		URL:            httppaths.GroupDownloader + httppaths.PathMediaItemShortLink,
		replaceInURL:   RowMenuActionItemIDKey,
		onlyStatusDone: true,
	},

	{RenderType: menuActionRenderTypeDivider},

	{
		RenderType:   menuActionRenderTypeAction,
		Action:       "delete",
		Title:        "Delete",
		icon:         icons.DownloaderRowMenuDeleteIcon,
		URL:          httppaths.GroupDownloader + httppaths.PathMediaItem,
		replaceInURL: RowMenuActionItemIDKey,
	},
}

func RowMenuActions(mapReplaceUrl map[string]string, isStatusDone bool) []rowMenuAction {
	actions := make([]rowMenuAction, 0, len(rowMenuActions))
	for _, a := range rowMenuActions {
		if !isStatusDone && a.onlyStatusDone {
			continue
		}

		actions = append(actions, a)
		action := &actions[len(actions)-1]

		if action.icon.FileName() != "" {
			action.IconSvg = action.icon.FileRaw()
		}

		if key := action.replaceInURL; key != "" {
			value, ok := mapReplaceUrl[key]
			if ok {
				action.URL = strings.Replace(action.URL, key, value, 1)
			}
		}
	}
	return actions
}
