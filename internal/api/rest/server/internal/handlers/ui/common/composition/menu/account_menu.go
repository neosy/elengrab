package menu

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type accountMenuAction struct {
	menuAction
	icon *icons.Icon

	requiredUserType dtypes.UserType
}

const (
	accountMenuActionMyDownloads = "my-downloads"
	accountMenuActionUsers       = "users"
	accountMenuActionLogout      = "logout"
)

var accountMenuActions = []accountMenuAction{
	{
		menuAction: menuAction{
			RenderType: renderTypeLink,
			Action:     accountMenuActionMyDownloads,
			Title:      "My Downloads",
			Link: linkOptions{
				NewTab: false,
			},
		},
		icon: &icons.UploadMenuMyDownloadsIcon,
	},
	{
		menuAction: menuAction{
			RenderType: renderTypeLink,
			Action:     accountMenuActionUsers,
			Title:      "Users",
			Link: linkOptions{
				URLTemplate: httppaths.AdminUsersPath,
				NewTab:      false,
			},
		},
		icon:             &icons.AdminUsersIcon,
		requiredUserType: dtypes.UserTypeAdmin,
	},
	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     accountMenuActionLogout,
			Title:      "Logout",
			Link: linkOptions{
				URLTemplate: httppaths.AuthLogoutPath,
				NewTab:      false,
			},
		},
		icon: &icons.AccountMenuLogoutIcon,
	},
}

func NewAccountMenuActions(authCtx dauth.AuthContext) []accountMenuAction {
	userType := authCtx.UserType()
	actions := make([]accountMenuAction, 0, len(accountMenuActions))

	for _, action := range accountMenuActions {
		if !userType.HasRequiredAccess(action.requiredUserType) {
			continue
		}

		switch action.Action {
		case accountMenuActionMyDownloads:
			if authCtx.Login == "" {
				continue
			}

			action.Link.URLTemplate = httppaths.BuildRootUserDownloadsPath(authCtx.Login)
		}

		if action.icon != nil {
			action.IconSvg = action.icon.FileRaw()
		}

		actions = append(actions, action)
	}
	return actions
}
