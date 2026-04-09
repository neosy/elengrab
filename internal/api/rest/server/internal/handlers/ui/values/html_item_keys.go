package uivalues

const (
	CssPathsKey        = "CssPaths"
	JsScriptsKey       = "JsScripts"
	JsImportMapJSONKey = "JsImportMapJSON"

	CssStyleKey = "CssStyle"

	DownloadStatusKey                     = "DownloadStatus"
	DownloadWorkingStatusKey              = "WorkingStatus"
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
	DownloadResultItemDeleteIconKey       = "GrabResultItemDeleteIcon"
	DownloadResultItemStatusFailedIconKey = "DownloadResultItemStatusFailedIcon"
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
	ViewerValues struct {
		ShowBackButton   bool
		IsVideoPlayer    bool
		MediaParametes   any
		MediaTitle       string
		MediaDescription string
		ContentType      string `json:"ContentType,omitempty"`
	}
)

func (v *ViewerValues) ToMap() map[string]any {
	return StructToMap(v)
}
