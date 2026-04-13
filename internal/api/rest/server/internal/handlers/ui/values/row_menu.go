package uivalues

import (
	"html/template"
	"strings"

	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

type rowMenuAction struct {
	Action  string `json:"action"`
	Title   string `json:"title"`
	IconSvg any    `json:"iconSvg,omitempty"`

	IconFileName   string
	URL            string
	NewTab         bool
	replaceInURL   string
	onlyStatusDone bool
}

const (
	RowMenuActionFileIdKey = "{fileId}"
	RowMenuActionURLKey    = "{url}"
)

var rowMenuActions = []rowMenuAction{
	{
		Action:         "watch",
		Title:          "Watch in new tab",
		IconFileName:   "menu-play-icon.svg",
		URL:            httppaths.GroupDownloader + httppaths.PathFileWatch,
		NewTab:         true,
		replaceInURL:   RowMenuActionFileIdKey,
		onlyStatusDone: true,
	},
	{
		Action:         "share-link",
		Title:          "Share link",
		IconFileName:   "menu-share-link-icon.svg",
		URL:            httppaths.GroupDownloader + httppaths.PathFileShortLink,
		replaceInURL:   RowMenuActionFileIdKey,
		onlyStatusDone: true,
	},
	{
		Action:       "open-original",
		Title:        "Open original in new tab",
		IconFileName: "menu-external-link-icon.svg",
		URL:          RowMenuActionURLKey,
		NewTab:       true,
		replaceInURL: RowMenuActionURLKey,
	},
	{
		Action:       "delete",
		Title:        "Delete",
		IconFileName: "download-delete-icon.svg",
		URL:          httppaths.GroupDownloader + httppaths.PathFile,
		replaceInURL: RowMenuActionFileIdKey,
	},
}

func RowMenuActions(svgDir string, mapReplaceUrl map[string]string, isStatusDone bool) []rowMenuAction {
	actions := make([]rowMenuAction, 0, len(rowMenuActions))
	for _, a := range rowMenuActions {
		if !isStatusDone && a.onlyStatusDone {
			continue
		}

		actions = append(actions, a)
		action := &actions[len(actions)-1]

		if action.IconFileName != "" {
			action.IconSvg = template.HTML(IconFileRaw(action.IconFileName, svgDir))
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
