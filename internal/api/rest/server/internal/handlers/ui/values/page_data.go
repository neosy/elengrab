package uivalues

import "html/template"

// Index page
type (
	IndexPageData struct {
		BasePaths  basePaths
		BaseValues baseValues
		Paths      PagePaths
		Values     IndexPageValues
		Extra      map[string]any
	}

	IndexPageValues struct {
		ShowHistorySearch bool

		DiskFree string
		DiskUsed string

		GrabForm IndexGrabForm
	}

	IndexGrabForm struct {
		InputPlaceholder string
	}
)

// Watch page
type (
	WatchPageData struct {
		BasePaths  basePaths
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

// Auth login page
type (
	AuthLoginPageData struct {
		BasePaths  basePaths
		BaseValues baseValues
		Paths      PagePaths
		Extra      map[string]any
	}
)

// Auth register page
type (
	AuthRegisterPageData struct {
		BasePaths  basePaths
		BaseValues baseValues
		Paths      PagePaths
		Extra      map[string]any
	}
)

// Error page
type (
	ErrorPageData struct {
		BasePaths  basePaths
		BaseValues baseValues
		Values     ErrorPageValues
		Extra      map[string]any
	}

	ErrorPageValues struct {
		Title   string
		Header  string
		BaseURL string

		CssStyle template.HTML

		ErrorCode  int
		ErrorTitle string
		ErrorText  string

		DebugErrorText any
		DebugData      template.HTML
	}
)

type (
	PagePaths struct {
		Css []string

		JsScripts       []jsScript
		JsImportMapJSON template.HTML

		PwaManifest string
		Stream      string
	}
)

// Page fragment
type (
	PageFragmentData struct {
		BasePaths  basePaths
		BaseValues baseValues
		Extra      map[string]any
	}
)

// Row fragment
type (
	RowFragmentData struct {
		BasePaths     basePaths
		BaseValues    baseValues
		Values        *RowFragmentValues
		IconFileNames map[string]any
		Extra         map[string]any
	}

	RowFragmentValues struct {
		DownloadID     string
		DownloadStatus string
		WorkingStatus  string

		DownloadRowPath    string
		DownloadRepeatPath string

		YoutubeChannelID string
		ThumbnailID      string
		AvatarTitle      string

		ImageURL       string
		ImageAvatarURL string
		ImageSiteURL   string

		MediaTitle string
		MediaURL   string

		ContentTimeAgo string

		FilePath      string
		FileSize      string
		Format        string
		DataFormat    string
		FormatTitle   string
		FormatTooltip string
		IsAudio       string
		DownloadURL   string
		StreamURL     string
		WatchURL      string
		DeleteURL     string
		RowID         string
		ProgressID    string

		DownloaderResultItemSourceLinkIcon   template.HTML
		DownloaderResultItemStatusIcon       template.HTML
		DownloaderResultItemDeleteIcon       template.HTML
		DownloaderResultItemStatusFailedIcon template.HTML
		IsItemHTMXOptionRepeat               bool
		PageHasDivItems                      bool
		ResultRowFade                        string
		ResultRowStatusTitle                 string
		ResultMediaUrlFade                   string
		ResultSizeFade                       string
		ResultFormatFade                     string
		IsDownloadEvent                      bool
		IsItemSpiner                         bool
	}
)
