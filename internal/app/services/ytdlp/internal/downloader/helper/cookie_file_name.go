package helper

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfile "github.com/neosy/elengrab/internal/pkg/filex"
)

func CookieFileNameFromURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", fmt.Errorf("The media URL is empty")
	}

	var cookieName string
	host := hostdetect.Detect(rawURL)
	if host != dtypes.MediaHostNone {
		cookieName = host.String()
	}

	if cookieName == "" {
		u, err := url.Parse(rawURL)
		if err != nil {
			return "", err
		}
		cookieName = strings.ToLower(u.Hostname())

	}

	return fmt.Sprintf("%s.txt", cookieName), nil
}

func CookieFilePath(fileName string, cookiesDir string) (string, error) {
	if fileName == "" {
		return "", fmt.Errorf("cookie file name is not specified")
	}

	if cookiesDir == "" {
		return "", fmt.Errorf("cookies dir is not specified")
	}

	cookieFilePath := filepath.Join(cookiesDir, fileName)

	// Check if the cookies file exists
	exists, err := nfile.FileExists(cookieFilePath)
	if err != nil {
		return "", fmt.Errorf("Failed check cookies file %s: %w", cookieFilePath, err)
	}
	if !exists {
		return "", fmt.Errorf("Cookies file '%s' not found", cookieFilePath)
	}
	return cookieFilePath, nil
}

func CookieFilePathFromURL(url string, cookiesDir string) (string, error) {
	fileName, err := CookieFileNameFromURL(url)
	if err != nil {
		return "", err
	}
	return CookieFilePath(fileName, cookiesDir)
}
