package handlers

import (
	"time"
)

const (
	dateFormate = time.RFC3339

	loadHistoryLimit    = 30
	preloadHistoryAfter = 29

	formFieldLoginKey           = "login"
	formFieldPasswordKey        = "password"
	formFieldConfirmPasswordKey = "confirmPassword"

	formFieldMediaURLKey          = "mediaURL"
	formFieldQualityCodecKey      = "quality-codec"
	formFieldQualityResolutionKey = "quality-resolution"
	formFieldFormatKey            = "format"

	urlKey  = "url"
	textKey = "text"

	sourceKey = "source"

	downloadIDKey        = "downloadID"
	beforeKey        = "before"
	channelIDKey     = "channelID"
	filterByTitleKey = "title"
	searchKey        = "search"
	shortCodeKey     = "shortCode"
)
