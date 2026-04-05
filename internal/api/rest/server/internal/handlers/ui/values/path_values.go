package uivalues

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

const (
	PathItemsHistoryKey = "PathItemsHistory"
)

var (
	PathValues = map[string]any{
		"PathStaticImg":   httppaths.GroupImg,
		"PathStaticIcons": httppaths.GroupIcon,
		"PathStaticCss":   httppaths.GroupCss,
		"PathStaticJs":    httppaths.GroupJs,
		"PathStaticPwa":   httppaths.GroupPwa,

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
