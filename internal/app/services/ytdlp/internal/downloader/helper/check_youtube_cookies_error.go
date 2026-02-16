package helper

import "strings"

func CheckForYouTubeCookiesError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "YouTube") &&
		strings.Contains(errStr, "Use --cookies-from-browser or --cookies")
}
