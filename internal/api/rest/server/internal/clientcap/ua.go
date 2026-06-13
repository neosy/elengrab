package clientcap

import (
	"regexp"
	"strconv"
)

var iosRegex = regexp.MustCompile(`OS (\d+)[_.](\d+)?[_.]?(\d+)?`)

func parseIOSVersion(ua string) (major, minor, patch int, ok bool) {
	m := iosRegex.FindStringSubmatch(ua)
	if len(m) == 0 {
		return 0, 0, 0, false
	}

	major = atoiSafe(m[1])
	minor = atoiSafe(m[2])
	patch = atoiSafe(m[3])

	return major, minor, patch, true
}

func atoiSafe(s string) int {
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}
