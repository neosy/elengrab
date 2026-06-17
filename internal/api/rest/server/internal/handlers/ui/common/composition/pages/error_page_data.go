package pages

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

// Error page
type (
	ErrorPageData struct {
		BasePaths  paths.HttpPaths
		BaseValues baseValues
		Values     ErrorPageValues
		Extra      map[string]any
	}

	ErrorPageValues struct {
		Title   string
		Header  string
		BaseURL string

		CssStyle template.HTML

		ErrorCode  int
		ErrorTitle string
		ErrorText  string

		DebugErrorText any
		DebugData      template.HTML
	}
)
