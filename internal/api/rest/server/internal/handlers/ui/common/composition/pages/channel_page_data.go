package pages

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

// Index page
type ChannelPageData struct {
	BasePaths  paths.HttpPaths
	BaseValues baseValues
	Paths      PagePaths
	Values     ChannelPageValues
	Extra      map[string]any
}

type ChannelPageValues struct {
	PagesListValues
}
