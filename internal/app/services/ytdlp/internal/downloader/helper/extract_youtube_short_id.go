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

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")

	// the link should look like:
	// /shorts/<id>
	if len(parts) < 2 || parts[0] != "shorts" {
		return "", fmt.Errorf("not a youtube short url")
	}

	shortID := parts[1]
	if shortID == "" {
		return "", fmt.Errorf("short id is emprty")
	}

	return shortID, nil
}
