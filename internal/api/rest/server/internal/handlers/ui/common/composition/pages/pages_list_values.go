package pages

import "html/template"

type PagesListValues struct {
	RowsFragmentValues

	HeaderActionsSearchButtonIcon   template.HTML
	HeaderActionsDownloadButtonIcon template.HTML
	SearchBackArrowIcon             template.HTML

	ShowHistorySearch        bool
	HeaderActionsAvatarTitle string

	SearchQuery string

	ChannelHeader ChannelHeader

	ActiveViewMode string
	ViewModeTabs   []ViewModeTab

	HasCreateAccess         bool
	HasWriteOperationAccess bool

	DiskFree string
	DiskUsed string

	VideoPreview VideoPreview
}

type VideoPreview struct {
	SoundOnIcon  template.HTML
	SoundOffIcon template.HTML
}
