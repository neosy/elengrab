package handlers

import (
	"time"
)

const (
	dateFormate = time.RFC3339

	loadHistoryLimit    = 20
	preloadHistoryAfter = 19

	formFieldLoginKey           = "login"
	formFieldPasswordKey        = "password"
	formFieldConfirmPasswordKey = "confirmPassword"

	formFieldMediaURLKey          = "mediaURL"
	formFieldQualityCodecKey      = "quality-codec"
	formFieldQualityResolutionKey = "quality-resolution"
	formFieldFormatKey            = "format"

	urlKey    = "url"
	textKey   = "text"
	sourceKey = "source"

	downloadIDKey    = "itemId"
	beforeKey        = "before"
	channelIDKey     = "channelId"
	filterByTitleKey = "title"
	searchKey        = "search"
	shortCodeKey     = "shortCode"
)
