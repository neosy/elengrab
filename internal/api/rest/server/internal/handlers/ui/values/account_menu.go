package uivalues

import (
	"html/template"

	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

type accountMenuAction struct {
	Action       string
	Title        string
	IconSvg      any
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
		actions[i].IconSvg = template.HTML(IconFileRaw(action.IconFileName, svgDir))
	}
	return actions
}
