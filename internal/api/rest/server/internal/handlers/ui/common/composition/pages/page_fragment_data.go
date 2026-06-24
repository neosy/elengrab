package pages

import "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"

// Page fragment
type (
	PageFragmentData struct {
		BasePaths  paths.HttpPaths
		BaseValues baseValues
		Extra      map[string]any
	}
)
