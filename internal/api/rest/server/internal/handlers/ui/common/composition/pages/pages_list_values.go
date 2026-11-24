package pages

import "html/template"

type PagesListValues struct {
	RowsFragmentValues

	HeaderActionsSearchButtonIcon   template.HTML
	HeaderActionsDownloadButtonIcon template.HTML
	HeaderActionsAvatarTitle        string

	UserAvatar UserAvatar

	ItemsCount int

	ShowHistorySearch   bool
	SearchBackArrowIcon template.HTML

	SearchQuery string

	UploadMenuContent template.HTML

	ActiveContentType string
	ContentTypeTabs   []ContentTypeTab

	ActiveViewMode string
	ViewModeTabs   []ViewModeTab

	HasCreateAccess         bool
	HasWriteOperationAccess bool

	VideoPreview VideoPreview

	AboutDialog        AboutDialogValues
	UploadFileDialog   UploadFileDialogValues
	DownloadListDialog DownloadListDialogValues
}

type UserAvatar struct {
	Icon       template.HTML
	ActionMode string
}

type VideoPreview struct {
	SoundOnIcon  template.HTML
	SoundOffIcon template.HTML
	PlayIcon     template.HTML
	PauseIcon    template.HTML
}
