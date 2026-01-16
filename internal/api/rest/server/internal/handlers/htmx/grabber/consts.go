package grabberh

import (
	"time"
)

const (
	dateFormate = time.RFC3339

	loadHistoryLimit = 30

	formFieldYouTubeURLKey        = "youtubeURL"
	formFieldQualityCodecKey      = "quality-codec"
	formFieldQualityResolutionKey = "quality-resolution"
	formFieldFormatKey            = "format"

	fileIdKey    = "fileId"
	beforeKey    = "before"
	channelIDKey = "channelID"
	userIDKey    = "userId"

	cookiePageHasDivItemsKey cookieKey = "page_has_div_items"

	channelIDValueNone = "none"
)
