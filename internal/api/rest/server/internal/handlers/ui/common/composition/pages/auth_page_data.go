package pages

import "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"

// Auth login page
type (
	AuthLoginPageData struct {
		BasePaths  paths.HttpPaths
		BaseValues baseValues
		Paths      PagePaths
		Values     AuthLoginPageValues
		Extra      map[string]any
	}

	AuthLoginPageValues struct {
		Redirect string
	}
)

// Auth register page
type (
	AuthRegisterPageData struct {
		BasePaths  paths.HttpPaths
		BaseValues baseValues
		Paths      PagePaths
		Extra      map[string]any
	}
)
