package uivalues

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

type accountMenuAction struct {
	Action       string `json:"action"`
	Title        string `json:"title"`
	IconSvg      any    `json:"iconSvg,omitempty"`
	IconFileName string
	URL          string
}

var accountMenuActions = []accountMenuAction{
	{
		Action:       "logout",
		Title:        "Logout",
		IconFileName: "menu-logout-icon.svg",
		URL:          httppaths.GroupAccount + httppaths.PathLogout,
	},
}

func AccountMenuActions(svgDir string) []accountMenuAction {
	actions := append([]accountMenuAction(nil), accountMenuActions...)
	for i, action := range actions {
		actions[i].IconSvg = IconFileRaw(action.IconFileName, svgDir)
	}
	return actions
}
