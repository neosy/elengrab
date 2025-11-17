package ytdownloader

import (
	"regexp"
	"strings"

	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

// sanitizeFileName replaces characters not allowed in file names with "_"
func (uc *YouTubeDownloader) sanitizeFileName(title string) string {
	// Trim spaces at the beginning and end
	title = strings.TrimSpace(title)

	// Replace all invalid characters with underscore
	// Windows forbidden: \ / : * ? " < > |
	// Unix forbidden: /
	safe := regexp.MustCompile(`[:|]`).ReplaceAllString(title, "-")
	safe = regexp.MustCompile(`[/\\]`).ReplaceAllString(safe, "-")
	safe = regexp.MustCompile(`["<>]`).ReplaceAllString(safe, "")
	safe = regexp.MustCompile(`[?*\x00-\x1F]`).ReplaceAllString(safe, "_")

	// Optional: replace multiple underscores with a single one
	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_")

	// Optional: limit length (e.g., 255 characters)
	if len(safe) > 255 {
		safe = safe[:255]
	}

	return safe
}

func (uc *YouTubeDownloader) enqueueDownloadTask(task *ddownload.DownloadTask) {
	uc.dlDispetcher.AddJob(wjobs.NewDownloadJob(task, uc))
}
