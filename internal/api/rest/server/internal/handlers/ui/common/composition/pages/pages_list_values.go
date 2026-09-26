package pages

import "html/template"

type PagesListValues struct {
	RowsFragmentValues

	HeaderActionsSearchButtonIcon   template.HTML
	HeaderActionsDownloadButtonIcon template.HTML
	HeaderActionsAvatarTitle        string

	UserAvatar UserAvatar

	ShowHistorySearch   bool
	SearchBackArrowIcon template.HTML

	SearchQuery string

	ActiveViewMode string
	ViewModeTabs   []ViewModeTab

	HasCreateAccess         bool
	HasWriteOperationAccess bool

	VideoPreview VideoPreview
}

type UserAvatar struct {
	Icon       template.HTML
	ActionMode string
}

type VideoPreview struct {
	SoundOnIcon  template.HTML
	SoundOffIcon template.HTML
}
