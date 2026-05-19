package menu

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

type accountMenuAction struct {
	Action  string `json:"action"`
	Title   string `json:"title"`
	IconSvg any    `json:"iconSvg,omitempty"`
	icon    icons.Icon
	URL     string
}

var accountMenuActions = []accountMenuAction{
	{
		Action: "logout",
		Title:  "Logout",
		icon:   icons.AccountMenuLogoutIcon,
		URL:    httppaths.GroupAccount + httppaths.PathLogout,
	},
}

func AccountMenuActions() []accountMenuAction {
	actions := append([]accountMenuAction(nil), accountMenuActions...)
	for i, action := range actions {
		actions[i].IconSvg = action.icon.FileRaw()
	}
	return actions
}
