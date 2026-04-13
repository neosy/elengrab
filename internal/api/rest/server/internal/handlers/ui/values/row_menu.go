package uivalues

import (
	"html/template"
	"strings"

	"github.com/google/uuid"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

type rowMenuAction struct {
	Action  string `json:"action"`
	Title   string `json:"title"`
	IconSvg any    `json:"iconSvg,omitempty"`

	IconFileName   string
	URL            string
	NewTab         bool
	isFileIDinURL  bool
	onlyStatusDone bool
}

var rowMenuActions = []rowMenuAction{
	{
		Action:         "watch",
		Title:          "Watch in new tab",
		IconFileName:   "menu-play-icon.svg",
		URL:            httppaths.GroupDownloader + httppaths.PathFileWatch,
		NewTab:         true,
		isFileIDinURL:  true,
		onlyStatusDone: true,
	},
	{
		Action:         "share-link",
		Title:          "Share link",
		IconFileName:   "menu-share-link-icon.svg",
		URL:            httppaths.GroupDownloader + httppaths.PathFileShortLink,
		isFileIDinURL:  true,
		onlyStatusDone: true,
	},
	{
		Action:        "delete",
		Title:         "Delete",
		IconFileName:  "download-delete-icon.svg",
		URL:           httppaths.GroupDownloader + httppaths.PathFile,
		isFileIDinURL: true,
	},
}

func RowMenuActions(svgDir string, fileID uuid.UUID, isStatusDone bool) []rowMenuAction {
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

		if action.isFileIDinURL {
			action.URL = strings.Replace(action.URL, "{fileId}", fileID.String(), 1)
		}
	}
	return actions
}
