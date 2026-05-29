package pages

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers/pages/page"
)

var (
	pages = []page.Page{
		{
			Name:       page.PageNameDashboard,
			Title:      page.PageNameDashboard.Title(),
			LogoSymbol: page.PageNameDashboard.LogoSymbol(),
			LogoIcon:   page.PageNameDashboard.LogoIcon(),

			ContentTemplateName: "admin-dashboard-content",
		},
		{
			Name:       page.PageNameUsers,
			Title:      page.PageNameUsers.Title(),
			LogoSymbol: page.PageNameUsers.LogoSymbol(),
			LogoIcon:   page.PageNameUsers.LogoIcon(),

			ContentTemplateName: "admin-users-content",
		},
	}

	pagesByName map[page.PageName]page.Page
)

func init() {
	pagesByName = make(map[page.PageName]page.Page, len(pages))

	for _, page := range pages {
		pagesByName[page.Name] = page
	}
}

func PageByName(name page.PageName) page.Page {
	return pagesByName[name]
}

func PageByURI(uri, prefix string) page.Page {
	name := page.ParseFromURI(uri, prefix)
	return pagesByName[name]
}
