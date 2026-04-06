package uivalues

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

const (
	PathItemsHistoryKey = "PathItemsHistory"
)

var (
	PathValues = map[string]any{
		"PathStaticImg":   httppaths.GroupStaticImg,
		"PathStaticIcons": httppaths.GroupStaticIcon,
		"PathStaticCss":   httppaths.GroupStaticCss,
		"PathStaticJs":    httppaths.GroupStaticJs,
		"PathStaticPwa":   httppaths.GroupStaticPwa,

		"PathAuthRegister": httppaths.GroupAccount + httppaths.PathRegister,
		"PathAuthLogin":    httppaths.GroupAccount + httppaths.PathLogin,

		"PathDownloader":     httppaths.GroupDownloader,
		"PathAccountMenu":    httppaths.GroupDownloader + httppaths.PathAccountMenu,
		"PathRowMenu":        httppaths.GroupDownloader + httppaths.PathFileMenu,
		PathItemsHistoryKey:  httppaths.GroupDownloader + httppaths.PathHistory,
		"PathDownloaderGrab": httppaths.GroupDownloader + httppaths.PathGrab,
		"PathHistorySearch":  httppaths.GroupDownloader + httppaths.PathSearch,
	}
)
