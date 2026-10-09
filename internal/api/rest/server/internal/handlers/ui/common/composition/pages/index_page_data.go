package pages

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

// Index page
type IndexPageData struct {
	BasePaths  paths.HttpPaths
	BaseValues baseValues
	Paths      PagePaths
	Values     IndexPageValues
	Extra      map[string]any
}

type IndexPageValues struct {
	PagesListValues
	PagesFooterValues

	GrabForm MediaGrabForm
}

type MediaGrabForm struct {
	InputPlaceholder   string
	GetButtonTitle     string
	SettingsButtonIcon template.HTML
	GetButtonIcon      template.HTML
}
