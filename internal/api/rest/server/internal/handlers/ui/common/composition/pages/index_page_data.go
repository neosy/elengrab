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
		RowsFragmentValues

		UserMenuSearchButtonIcon   template.HTML
		UserMenuDownloadButtonIcon template.HTML
		SearchBackArrowIcon        template.HTML

		ShowHistorySearch   bool
		UserMenuAvatarTitle string

		SearchQuery string

		ChannelHeader ChannelHeader

		ActiveViewMode string
		ViewModeTabs   []ViewModeTab

		HasCreateAccess         bool
		HasWriteOperationAccess bool

		DiskFree string
		DiskUsed string

		GrabForm IndexGrabForm

		VideoPreview VideoPreview
	}

	IndexGrabForm struct {
		InputPlaceholder   string
		GetButtonTitle     string
		SettingsButtonIcon template.HTML
		GetButtonIcon      template.HTML
	}

	VideoPreview struct {
		SoundOnIcon  template.HTML
		SoundOffIcon template.HTML
	}
)
