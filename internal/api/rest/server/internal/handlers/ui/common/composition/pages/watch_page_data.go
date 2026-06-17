package pages

import "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"

// Watch page
type (
	WatchPageData struct {
		BasePaths  paths.HttpPaths
		BaseValues baseValues
		Paths      PagePaths
		Values     WatchPageValues
		Extra      map[string]any
	}

	WatchPageValues struct {
		IsVideoPlayer bool

		ShowBackButton bool
		ShowControls   bool
		ShowPlaylist   bool

		ContentType string `json:"ContentType,omitempty"`

		MediaTitle         string
		MediaTitleImageURL string
		MediaDescription   string
		MediaParameters    []MediaParameter
	}

	MediaParameter struct {
		Name  string
		Value string
		URL   string
	}
)
