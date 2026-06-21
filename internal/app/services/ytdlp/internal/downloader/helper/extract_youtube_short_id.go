package helper

import (
	"fmt"
	"net/url"
	"strings"

	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
)

func ExtractYouTubeShortID(rawURL string) (string, error) {
	if !hostdetect.YouTube(rawURL) {
		return "", fmt.Errorf("not a youtube url")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// the link should look like:
	// /shorts/<id>

	path := strings.Trim(u.Path, "/")

	prefix, shortID, ok := strings.Cut(path, "/")
	if !ok || prefix != "shorts" {
		return "", fmt.Errorf("not a youtube short url")
	}

	if shortID == "" {
		return "", fmt.Errorf("youtube short id not found")
	}

	if !youtubeIDRegex.MatchString(shortID) {
		return "", fmt.Errorf("invalid short id")
	}

	return shortID, nil
}
