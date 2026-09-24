package qkeys

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// keys holds the internal concrete implementation of the query keys.
var keys = NewQueryKeys()

// Keys is a globally accessible, interface-protected instance for key operations.
var Keys QueryKeysRegistry = keys

var (
	UrlKey    = keys.add("url", "u")
	TextKey   = keys.add("text", "t")
	SourceKey = keys.add("source", "src")

	DownloadIDKey = keys.add("itemId", "id")

	ChannelIDKey       = keys.addWithName("channelId", "cid", dtypes.QueryFilterNameChannelID)
	ChannelPlatformKey = keys.add("platform", "cpl")

	SearchKey      = keys.addWithName("query", "qr", dtypes.QueryFilterNameSearchQuery)
	SearchQueryKey = keys.addWithName("searchQuery", "sq", dtypes.QueryFilterNameSearchQuery)

	ShortCodeKey = keys.add("shortCode", "sc")
	RedirectKey  = keys.add("redirect", "r")

	ViewModeKey = keys.add("viewMode", "vm")

	SearchParametersKey = keys.add("searchParameters", "sp")

	LastIDKey       = keys.add("lastId", "lid")
	LastCreateAtKey = keys.add("lastCreateAt", "date")
	LastViewsKey    = keys.add("lastViews", "views")
	LastCursorKey   = keys.add("lastCursor", "lc")
)
