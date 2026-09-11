package qkeys

var Keys = NewQueryKeys()

var (
	UrlKey    = Keys.add("url", "u")
	TextKey   = Keys.add("text", "t")
	SourceKey = Keys.add("source", "src")

	DownloadIDKey       = Keys.add("itemId", "id")
	ChannelIDKey        = Keys.add("channelId", "cid")
	SearchKey           = Keys.add("search", "s")
	ShortCodeKey        = Keys.add("shortCode", "sc")
	RedirectKey         = Keys.add("redirect", "r")
	ViewModeKey         = Keys.add("viewMode", "vm")
	LastIDKey           = Keys.add("lastId", "lid")
	LastCreateAtKey     = Keys.add("lastCreateAt", "date")
	LastViewsKey        = Keys.add("lastViews", "views")
	SearchParametersKey = Keys.add("sp", "sp")
	LastCursorKey       = Keys.add("lastCursor", "lc")
)
