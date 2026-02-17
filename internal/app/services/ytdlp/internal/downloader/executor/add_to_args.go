package executor

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/pkg/nfile"
)

func addYouTubeCookiesToArgs(logger *slog.Logger, args []string, serviceOptions *dto.Options) []string {
	// Check if cookies are allowed in the service options
	if serviceOptions.YoutubeAllowCookies {
		path, err := ensureYouTubeCookiePath(serviceOptions.CookiesDir)
		if err != nil {
			logger.Warn("Failed YouTube cookie file", "error", err)
		}
		if path != "" {
			// Append the cookies file path to the arguments if it exists
			args = append(args, "--cookies", path)
		}
	}
	return args
}

func ensureYouTubeCookiePath(cookiesDir string) (string, error) {
	cookieFilePath := filepath.Join(cookiesDir, consts.YtDlpYouTubeCookieFileName)
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
