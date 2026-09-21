package pages

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

// Row fragment
type (
	RowsFragmentData struct {
		BasePaths  paths.HttpPaths
		BaseValues baseValues
		Values     *RowsFragmentValues
		Extra      map[string]any
	}

	RowsFragmentValues struct {
		HasSearchFilters     bool
		SearchParametersJSON string
		SearchParameters     template.HTML

		ChannelJSON string

		ResultNoRows   bool
		ResultRowsHTML template.HTML
	}
)
