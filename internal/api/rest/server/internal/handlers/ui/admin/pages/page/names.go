package page

import (
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
)

type PageName string

const (
	PageNameDashboard PageName = "dashboard"
	PageNameUsers     PageName = "users"
	PageNameGroups    PageName = "groups"
	PageNameLogs      PageName = "logs"
	PageNameSettings  PageName = "settings"

	pageNameDefault = PageNameDashboard
)

var (
	validPageNames = map[PageName]struct{}{
		PageNameDashboard: {},
		PageNameUsers:     {},
		PageNameGroups:    {},
		PageNameLogs:      {},
		PageNameSettings:  {},
	}

	titlesByName = map[PageName]string{
		PageNameDashboard: "Dashboard",
		PageNameUsers:     "Users",
		PageNameGroups:    "Groups",
		PageNameLogs:      "Activity Logs",
		PageNameSettings:  "Settings",
	}

	logoSymbolsByName = map[PageName]string{
		PageNameDashboard: "📊",
		PageNameUsers:     "👥",
		PageNameGroups:    "🔑",
		PageNameLogs:      "📋",
		PageNameSettings:  "⚙️",
	}

	logoIconByName = map[PageName]icons.Icon{
		PageNameUsers: icons.AdminUsersIcon,
	}
)

func (name PageName) String() string {
	return string(name)
}

func (name PageName) Title() string {
	return titlesByName[name]
}

func (name PageName) LogoSymbol() string {
	return logoSymbolsByName[name]
}

func (name PageName) LogoIcon() icons.Icon {
	return logoIconByName[name]
}

func ParseFromURI(uri, prefix string) PageName {
	pagePath := strings.TrimPrefix(uri, prefix)
	pagePath = strings.TrimPrefix(pagePath, "/")
	pagePathParts := strings.Split(pagePath, "/")

	var name string
	if len(pagePathParts) > 0 {
		name = pagePathParts[0]
	}

	_, exists := validPageNames[PageName(name)]
	if exists {
		return PageName(name)
	}

	return pageNameDefault
}
