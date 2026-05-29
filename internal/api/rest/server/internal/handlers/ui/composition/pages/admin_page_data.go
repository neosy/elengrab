package pages

import (
	"html/template"

	navmenu "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers/nav_menu"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
)

// Admin page
type (
	AdminPageData struct {
		BasePaths  paths.Paths
		BaseValues baseValues
		Paths      PagePaths
		Values     AdminPageValues
		Extra      map[string]any
	}

	AdminPageValues struct {
		PageTitle        string
		PageName         string
		IsPageLogoSymbol bool
		PageLogoHTML     template.HTML
		NavMenu          []navmenu.MenuItem

		ContentTemplateName string
		ContentHTML         template.HTML
	}
)
