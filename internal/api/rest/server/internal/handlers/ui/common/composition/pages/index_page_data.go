package pages

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

// Index page
type (
	IndexPageData struct {
		BasePaths  paths.HttpPaths
		BaseValues baseValues
		Paths      PagePaths
		Values     IndexPageValues
		Extra      map[string]any
	}

	IndexPageValues struct {
		UserMenuSearchButtonIcon   template.HTML
		UserMenuDownloadButtonIcon template.HTML
		SearchBackArrowIcon        template.HTML

		ShowHistorySearch   bool
		UserMenuAvatarTitle string

		HasCreateAccess bool
		HasWriteAccess  bool

		DiskFree string
		DiskUsed string

		GrabForm IndexGrabForm
	}

	IndexGrabForm struct {
		InputPlaceholder   string
		GetButtonTitle     string
		SettingsButtonIcon template.HTML
		GetButtonIcon      template.HTML
	}
)
