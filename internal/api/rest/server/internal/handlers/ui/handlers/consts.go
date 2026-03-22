package handlers

import (
	"time"
)

const (
	dateFormate = time.RFC3339

	loadHistoryLimit    = 30
	preloadHistoryAfter = 20

	formFieldMediaURLKey          = "mediaURL"
	formFieldQualityCodecKey      = "quality-codec"
	formFieldQualityResolutionKey = "quality-resolution"
	formFieldFormatKey            = "format"

	fileIdKey        = "fileId"
	beforeKey        = "before"
	channelIDKey     = "channelID"
	userIDKey        = "userId"
	filterByTitleKey = "title"
	searchKey        = "search"

	cookiePageHasDivItemsKey cookieKey = "page_has_div_items"

	channelIDValueNone = "none"
)
