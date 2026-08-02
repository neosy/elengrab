package pages

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

// Row fragment
type (
	RowFragmentData struct {
		BasePaths     paths.HttpPaths
		BaseValues    baseValues
		Values        *RowFragmentValues
		IconFileNames map[string]string
		Extra         map[string]any
	}

	RowFragmentValues struct {
		DownloadID     string
		DownloadStatus string
		WorkingStatus  string
		Visibility     string
		IsReady        bool

		DownloadRowPath    string
		DownloadRepeatPath string

		YoutubeChannelID string
		AvatarTitle      string

		ThumbnailID         string
		ThumbnailIsPortrait bool
		ThumbnailURL        string

		ImageURL       string
		ImageAvatarURL string
		ImageSiteURL   string

		MediaTitle string
		MediaURL   string

		Duration         string
		ContentViewCount string
		ContentTimeAgo   string

		WatchPercent float64
		Watched      bool

		UserName string

		FilePath string
		FileSize string

		Format        string
		DataFormat    string
		FormatTitle   string
		FormatTooltip string

		IsAudio      string
		VideoIsShort bool

		DownloadURL string
		StreamURL   string
		WatchURL    string
		DeleteURL   string

		ItemID     string
		RowID      string
		ProgressID string

		DownloaderResultItemSourceLinkIcon   template.HTML
		DownloaderResultItemStatusIcon       template.HTML
		DownloaderResultItemDeleteIcon       template.HTML
		DownloaderResultItemStatusFailedIcon template.HTML
		RefreshingIcon                       template.HTML
		MetaUserNameSeparatorIcon            template.HTML
		PublicIcon                           template.HTML
		PrivateIcon                          template.HTML
		ShareLinkIcon                        template.HTML
		WatchedIcon                          template.HTML

		IsItemHTMXOptionRepeat bool
		PageHasDivItems        bool
		ResultRowFade          string
		ResultRowStatusTitle   string
		ResultMediaUrlFade     string
		ResultSizeFade         string
		ResultFormatFade       string
		IsDownloadEvent        bool
		IsItemSpiner           bool
	}
)
