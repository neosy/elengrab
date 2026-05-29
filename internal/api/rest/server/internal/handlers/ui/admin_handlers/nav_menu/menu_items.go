package navmenu

import (
	"slices"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers/pages/page"
)

var (
	menuItems = []MenuItem{
		{
			PageName:   page.PageNameDashboard,
			Title:      page.PageNameDashboard.Title(),
			logoSymbol: page.PageNameDashboard.LogoSymbol(),
			logoIcon:   page.PageNameDashboard.LogoIcon(),

			IsVisible: true,
		},
		{
			PageName:   page.PageNameUsers,
			Title:      page.PageNameUsers.Title(),
			logoSymbol: page.PageNameUsers.LogoSymbol(),
			logoIcon:   page.PageNameUsers.LogoIcon(),

			IsVisible: true,
		},
		{
			PageName:   page.PageNameGroups,
			Title:      page.PageNameGroups.Title(),
			logoSymbol: page.PageNameGroups.LogoSymbol(),
			logoIcon:   page.PageNameGroups.LogoIcon(),

			IsVisible: false,
		},
		{
			PageName:   page.PageNameLogs,
			Title:      page.PageNameLogs.Title(),
			logoSymbol: page.PageNameLogs.LogoSymbol(),
			logoIcon:   page.PageNameLogs.LogoIcon(),

			IsVisible: false,
		},
		{
			PageName:   page.PageNameSettings,
			Title:      page.PageNameSettings.Title(),
			logoSymbol: page.PageNameSettings.LogoSymbol(),
			logoIcon:   page.PageNameSettings.LogoIcon(),

			IsVisible: false,
		},
	}

	menuItemByPageName map[page.PageName]MenuItem
)

func init() {
	menuItemByPageName = make(map[page.PageName]MenuItem, len(menuItems))

	for i, item := range menuItems {
		item.Order = uint8(i)
		item.URL = item.url()

		menuItems[i] = item
		menuItemByPageName[item.PageName] = item
	}
}

func NewMenuItems(activePage page.PageName) []MenuItem {
	items := slices.Clone(menuItems)

	activeItem := menuItemByPageName[activePage]
	items[activeItem.Order].IsActive = true

	return items
}
