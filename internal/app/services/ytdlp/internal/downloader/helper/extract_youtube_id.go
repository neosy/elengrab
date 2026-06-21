package helper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
)

var youtubeIDRegex = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

func ExtractYouTubeID(rawURL string) (string, error) {
	// https://youtu.be/Ulpa7scJnPw?si=4BAxG3K-p_FyhmDJ
	if !hostdetect.YouTube(rawURL) {
		return "", fmt.Errorf("not a youtube url")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	var youtubeID string

	// the link should look like:
	// /<id>
	// /watch?v=<id>

	switch strings.Trim(u.Path, "/") {
	case "watch":
		youtubeID = u.Query().Get("v")
	default:
		youtubeID = strings.Trim(u.Path, "/")
	}

	if youtubeID == "" {
		return "", fmt.Errorf("youtube id not found")
	}

	if !youtubeIDRegex.MatchString(youtubeID) {
		return "", fmt.Errorf("invalid youtube id")
	}

	return youtubeID, nil
}
