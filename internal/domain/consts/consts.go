package dconsts

const (
	MediaTitleMaxLength       = 100
	MediaDescriptionMaxLength = 5000

	// MediaURLMaxLength is the maximum length of a single media or list URL.
	MediaURLMaxLength = 2048

	// MediaDownloadListMaxLen is the maximum number of media that can be added to downloads in one request.
	// The extract endpoint is capped by the same value, so a user never selects more than the server accepts.
	MediaDownloadListMaxLen = 5000

	// DefaultIncludeTitleInFilename is the default for the "include title in file name" option.
	DefaultIncludeTitleInFilename = true

	// DefaultDownloadChannelImage is the default for the "download channel image" option.
	DefaultDownloadChannelImage = true
)
