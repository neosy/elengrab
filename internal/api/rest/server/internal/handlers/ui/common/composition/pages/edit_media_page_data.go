package pages

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

// Edit media page
type (
	EditMediaPageData struct {
		BasePaths  paths.HttpPaths
		BaseValues baseValues
		Paths      PagePaths
		Values     EditMediaPageValues
		Extra      map[string]any
	}

	EditMediaPageValues struct {
		MediaTitle       string
		MediaDescription string

		PatchURL string

		Visibility     string
		VisibilityList []VisibilityOption
	}
)
