package youtube

import "regexp"

var youtubeIDRegex = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)
