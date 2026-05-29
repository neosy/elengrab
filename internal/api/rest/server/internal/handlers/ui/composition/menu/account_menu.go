package menu

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type accountMenuAction struct {
	Action         string        `json:"action"`
	Title          string        `json:"title"`
	IconSvg        template.HTML `json:"iconSvg,omitempty"`
	icon           icons.Icon
	URL            string
	accessUserType dtypes.UserType
}

var accountMenuActions = []accountMenuAction{
	{
		Action:         "users",
		Title:          "Users",
		icon:           icons.AdminUsersIcon,
		URL:            httppaths.AdminGroup + httppaths.AdminUsersPath,
		accessUserType: dtypes.UserTypeAdmin,
	},
	{
		Action: "logout",
		Title:  "Logout",
		icon:   icons.AccountMenuLogoutIcon,
		URL:    httppaths.GroupAccount + httppaths.PathLogout,
	},
}

func NewAccountMenuActions(user *dauth.UserContext) []accountMenuAction {
	actions := make([]accountMenuAction, 0, len(accountMenuActions))
	for _, action := range accountMenuActions {
		if action.accessUserType > user.UserType() {
			continue
		}

		action.IconSvg = action.icon.FileRaw()

		actions = append(actions, action)
	}
	return actions
}
