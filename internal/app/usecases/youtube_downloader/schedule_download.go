package ytdownloader

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (uc *YouTubeDownloader) ScheduleDownload(
	ctx context.Context,
	url string,
	options *ddownload.DownloadOptions,
) (*dto.ScheduleDownloadResponse, error) {
	fileId := uuid.New()
	filename := fileId.String()

	options.Filename = &filename

	err := uc.file.Create(
		ctx,
		&ddownload.File{
			FileId:     fileId,
			FileName:   filename,
			YoutubeUrl: url,
		},
		options,
	)
	if err != nil {
		uc.logger.Error("Insert record failed", "error", err)
		return nil, err
	}

	file, err := uc.file.FindByFileId(ctx, fileId, true)
	if err != nil {
		return nil, err
	}

	err = uc.fileStatus.Pending(ctx, fileId)
	if err != nil {
		uc.fileStatus.Failed(ctx, fileId, uptr.String(err.Error()))
		return nil, err
	}

	uc.enqueueDownloadTask(file.DownloadTask)

	return &dto.ScheduleDownloadResponse{
		FileId:       file.FileId,
		Status:       file.Status,
		YoutubeTitle: file.YoutubeTitle,
		Format:       file.Ext,
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

func (uc *YouTubeDownloader) enqueueDownloadTask(task *ddownload.DownloadTask) {
	uc.dlDispetcher.AddJob(wjobs.NewDownloadJob(task, uc))
}
