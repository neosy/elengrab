package qkeys

import dtypes "github.com/neosy/elengrab/internal/domain/types"

var Keys = newQueryKeys()

var (
	UrlKey    = Keys.add("url", "u")
	TextKey   = Keys.add("text", "t")
	SourceKey = Keys.add("source", "src")

	DownloadIDKey       = Keys.add("itemId", "id")
	ChannelIDKey        = Keys.addWithName("channelId", "cid", dtypes.QueryFilterNameChannelID)
	ChannelPlatformKey  = Keys.add("platform", "cpl")
	SearchKey           = Keys.addWithName("query", "qr", dtypes.QueryFilterNameSearchQuery)
	SearchQueryKey      = Keys.addWithName("searchQuery", "sq", dtypes.QueryFilterNameSearchQuery)
	ShortCodeKey        = Keys.add("shortCode", "sc")
	RedirectKey         = Keys.add("redirect", "r")
	ViewModeKey         = Keys.add("viewMode", "vm")
	LastIDKey           = Keys.add("lastId", "lid")
	LastCreateAtKey     = Keys.add("lastCreateAt", "date")
	LastViewsKey        = Keys.add("lastViews", "views")
	SearchParametersKey = Keys.add("searchParameters", "sp")
	LastCursorKey       = Keys.add("lastCursor", "lc")
)
