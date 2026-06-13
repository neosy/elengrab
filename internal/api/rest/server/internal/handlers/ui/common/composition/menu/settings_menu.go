package menu

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

type settingsMenuAction struct {
	menuAction
	icon *icons.Icon
}

var settingsMenuActions = []settingsMenuAction{
	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     "grid-view",
			Title:      "Grid view: on",
			Text:       "Grid view",
			Link: linkOptions{
				URL: httppaths.GroupDownloader + httppaths.PathSettingsMenu,
			},
		},
		icon: nil,
	},
}

func SettingsMenuActions() []settingsMenuAction {
	actions := append([]settingsMenuAction(nil), settingsMenuActions...)
	for _, action := range actions {
		if action.icon != nil {
			action.IconSvg = action.icon.FileRaw()
		}
	}
	return actions
}
