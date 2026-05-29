package navmenu

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers/pages/page"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

type MenuItem struct {
	Title      string
	PageName   page.PageName
	logoSymbol string
	logoIcon   icons.Icon
	URL        string

	Order uint8

	IsActive  bool
	IsVisible bool
}

func (m *MenuItem) HasIcon() bool {
	return !m.logoIcon.IsZero()
}

func (m *MenuItem) LogoSymbolHTML() template.HTML {
	return template.HTML(m.logoSymbol)
}

func (m *MenuItem) LogoIconHTML() template.HTML {
	return m.logoIcon.FileRaw()
}

func (m *MenuItem) LogoHTML() template.HTML {
	if m.HasIcon() {
		return m.LogoIconHTML()
	}
	return m.LogoSymbolHTML()
}

func (mi MenuItem) url() string {
	path := httppaths.AdminGroup

	switch mi.PageName {
	case page.PageNameUsers:
		path += httppaths.AdminUsersPath
	case page.PageNameGroups:
		path += httppaths.AdminGroupsPath
	case page.PageNameLogs:
		path += httppaths.AdminLogsPath
	case page.PageNameSettings:
		path += httppaths.AdminSettingsPath
	}

	return path
}
