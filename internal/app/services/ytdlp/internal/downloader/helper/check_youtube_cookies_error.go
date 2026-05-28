package helper

import "strings"

func CheckCookiesError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "Use --cookies-from-browser or --cookies")
}
