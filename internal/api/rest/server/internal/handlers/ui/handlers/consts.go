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
	userKey          = "user"
	filterByTitleKey = "title"
	searchKey        = "search"

	cookiePageHasDivItemsKey cookieKey = "page_has_div_items"

	channelIDValueNone = "none"

	dirStaticName = "static"
	dirCssName    = "css"
	dirJsName     = "js"
)
