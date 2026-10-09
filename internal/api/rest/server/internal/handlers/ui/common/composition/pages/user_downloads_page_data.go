package pages

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

// User downloads page
type UserDownloadsPageData struct {
	BasePaths  paths.HttpPaths
	BaseValues baseValues
	Paths      PagePaths
	Values     UserDownloadsPageValues
	Extra      map[string]any
}

type UserDownloadsPageValues struct {
	PagesListValues
	PagesFooterValues

	GrabForm MediaGrabForm
	UserHeader UserDownloadsHeader
}

type UserDownloadsHeader struct {
	ListHeader
}
