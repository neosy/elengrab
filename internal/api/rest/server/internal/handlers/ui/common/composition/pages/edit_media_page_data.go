package pages

import (
	"slices"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
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
		VisibilityList []EditMediaPageVisibility
	}

	EditMediaPageVisibility struct {
		Value string
		Label string
	}
)

var (
	editMediaPageVisibilityList []EditMediaPageVisibility
)

func init() {
	editMediaPageVisibilityList = newMediaPageVisibilityList()
}

func MediaPageVisibilityList() []EditMediaPageVisibility {
	return slices.Clone(editMediaPageVisibilityList)
}

func newMediaPageVisibilityList() []EditMediaPageVisibility {
	var list = make([]EditMediaPageVisibility, len(dtypes.MediaVisibilityList()))
	for i, v := range dtypes.MediaVisibilityList() {
		list[i] = EditMediaPageVisibility{
			Value: v.String(),
			Label: v.Label(),
		}
	}
	return list
}
