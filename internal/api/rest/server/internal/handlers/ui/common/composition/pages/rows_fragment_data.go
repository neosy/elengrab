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
		SearchFiltersJSON    string
		SearchParametersJSON string
		SearchParameters     template.HTML

		ResultNoRows   bool
		ResultRowsHTML template.HTML

		VideoPlayerMinimizeIcon template.HTML

		PlayIcon  template.HTML
		PauseIcon template.HTML

		VideoPlayingIcon template.HTML
		AudioPlayingIcon template.HTML
		AudioNoteIcon    template.HTML
	}
)
