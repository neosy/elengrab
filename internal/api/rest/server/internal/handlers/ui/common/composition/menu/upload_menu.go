package menu

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

type uploadMenuAction struct {
	menuAction
	icon *icons.Icon
}

const (
	uploadMenuActionUploadURL   = "upload-url"
	uploadMenuActionMyDownloads = "my-downloads"
)

var uploadMenuActions = []uploadMenuAction{
	{
		menuAction: menuAction{
			RenderType: renderTypeAction,
			Action:     uploadMenuActionUploadURL,
			Title:      "Add download",
		},
		icon: &icons.UploadMenuPlusIcon,
	},
	{
		menuAction: menuAction{
			RenderType: renderTypeLink,
			Action:     uploadMenuActionMyDownloads,
			Title:      "My Downloads",
			Link: linkOptions{
				URL:    RowMenuActionURLKey,
				NewTab: false,
			},
		},
		icon: &icons.UploadMenuMyDownloadsIcon,
	},
}

func UploadMenuActions(userLogin string) []uploadMenuAction {
	actions := make([]uploadMenuAction, 0, len(uploadMenuActions))

	for _, action := range uploadMenuActions {
		switch action.Action {
		case uploadMenuActionMyDownloads:
			if userLogin == "" {
				continue
			}

			action.Link.URL = httppaths.BuildRootUserDownloadsPath(userLogin)
		}

		if action.icon != nil {
			action.IconSvg = action.icon.FileRaw()
		}

		actions = append(actions, action)
	}
	return actions
}
