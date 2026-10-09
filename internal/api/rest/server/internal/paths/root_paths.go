package httppaths

import "strings"

// Root
const (
	IndexPath = "/"

	RootFaviconICOPath = "/favicon.ico"
	RootRobotsTxtPath  = "/robots.txt"

	RootUserLoginPath = "/@{login}"
)

func buildRootUserLoginPath(path string, login string) string {
	return strings.Replace(path, "{login}", login, 1)
}

func BuildRootUserDownloadsPath(login string) string {
	return buildRootUserLoginPath(RootUserLoginPath, login)
}
