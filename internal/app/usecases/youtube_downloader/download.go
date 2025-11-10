package ucdownloader

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *YouTubeDownloader) Download(
	ctx context.Context,
	url string,
	options *ddownload.DownloadOptions,
) (*dto.DownloadResponse, error) {
	fileId := uuid.New()
	filename := fileId.String()

	options.Filename = &filename

	err := uc.fileRep.Insert(
		ctx,
		&ddownload.File{
			FileId:   fileId,
			FileName: filename,
		},
	)
	if err != nil {
		uc.logger.Error("Insert record failed", "error", err)
	}

	result, err := uc.downloaderSrv.Download(
		url,
		options,
	)

	if err != nil {
		return nil, err
	}

	safeReadableFullName := fmt.Sprintf("%s.%s", uc.sanitizeFileName(result.Title), result.FileExt)

	err = uc.fileRep.Update(
		ctx,
		&ddownload.File{
			FileId:               fileId,
			Title:                result.Title,
			FileName:             result.Filename,
			Ext:                  result.FileExt,
			FullName:             result.FileFullName,
			SafeReadableFullName: safeReadableFullName,
		})
	if err != nil {
		uc.logger.Error("Update record error", "error", err)
	}

	return &dto.DownloadResponse{
		Title:  result.Title,
		Format: result.FileExt,
		FileId: fileId,
	}, nil
}

// sanitizeFileName replaces characters not allowed in file names with "_"
func (uc *YouTubeDownloader) sanitizeFileName(title string) string {
	// Trim spaces at the beginning and end
	title = strings.TrimSpace(title)

	// Replace all invalid characters with underscore
	// Windows forbidden: \ / : * ? " < > |
	// Unix forbidden: /
	safe := regexp.MustCompile(`[:|]`).ReplaceAllString(title, "-")
	safe = regexp.MustCompile(`["<>]`).ReplaceAllString(safe, "")
	safe = regexp.MustCompile(`[/\\?*\x00-\x1F]`).ReplaceAllString(safe, "_")

	// Optional: replace multiple underscores with a single one
	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_")

	// Optional: limit length (e.g., 255 characters)
	if len(safe) > 255 {
		safe = safe[:255]
	}

	return safe
}
