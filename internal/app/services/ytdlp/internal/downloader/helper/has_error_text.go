package helper

import "strings"

func HasErrTextForbidden(err error) bool {
	return err != nil && strings.Contains(err.Error(), "HTTP Error 403")
}
