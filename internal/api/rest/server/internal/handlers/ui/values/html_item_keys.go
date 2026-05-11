package uivalues

const (
	CssPathsKey        = "CssPaths"
	JsScriptsKey       = "JsScripts"
	JsImportMapJSONKey = "JsImportMapJSON"
	PwaPathsKey        = "PwaPaths"
	PwaManifestPathKey = "PwaManifestPath"

	CssStyleKey = "CssStyle"

	DownloadingProgressPercentKey         = "ProgressPercent"
	ResultNoRowsKey                       = "ResultNoRows"
	ResultRowsHTMLKey                     = "ResultRowsHTML"
	ResultRowStatusIconKey                = "ResultRowStatusIcon"
	ResultRowStatusTitleKey               = "ResultRowStatusTitle"
	ResultRowFadeKey                      = "ResultRowFade"
	PageHasDivItemsKey                    = "PageHasDivItems"
	IsItemHTMXOptionRepeatKey             = "IsItemHTMXOptionRepeat"
	IsFileEventKey                        = "IsFileEvent"
	IsItemSpinerKey                       = "IsItemSpiner"
	DownloaderResultItemDeleteIconKey       = "GrabResultItemDeleteIcon"
	DownloaderResultItemStatusFailedIconKey = "DownloaderResultItemStatusFailedIcon"
	ShowHistorySearchKey                  = "ShowHistorySearch"
	PathStreamKey                         = "PathStream"

	AuthLoginKey    = "Login"
	AuthErrorMsgKey = "ErrorMsg"

	UserAvatarIconKey       = "UserAvatarIcon"
	UserAvatarActionModeKey = "UserAvatarActionMode"
	UserLoginKey            = "UserLogin"
	UserEmailKey            = "UserEmail"

	AccountMenuActionsKey = "AccountMenuActions"
	RowMenuActionsKey     = "RowMenuActions"

	AppVersionKey = "AppVersion"
	DiskFreeKey   = "DiskFree"
	DiskUsedKey   = "DiskUsed"

	ResultMediaUrlFadeKey = "ResultMediaUrlFade"
	ResultSizeFadeKey     = "ResultSizeFade"
	ResultFormatFadeKey   = "ResultFormatFade"
	DisableHTMXEventKey   = "DisableHTMXEvent"

	DebugDataKey = "DebugData"
)

type (
	WatcherValues struct {
		ShowBackButton   bool
		IsVideoPlayer    bool
		MediaParametes   any
		MediaTitle       string
		MediaDescription string
		ContentType      string `json:"ContentType,omitempty"`
	}
)

func (v *WatcherValues) ToMap() map[string]any {
	return StructToMap(v)
}
