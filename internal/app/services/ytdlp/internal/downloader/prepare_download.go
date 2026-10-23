package downloader

import (
	"context"
	"fmt"

	"uuid"

	downloadpreparer "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/download_preparer"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/utils/filename"
	"github.com/neosy/elengrab/internal/pkg/filex"
)

func (d *Downloader) prepareDownload(
	ctx context.Context,
	rawURL string,
	options idto.DLOptions,
) (*idto.DownloadMeta, *idto.DownloadExecOptions, error) {
	// Build yt-dlp arguments and get file extension and title
	preparer := downloadpreparer.NewDownloadPreparer(d.executor.FetchInfoWithBestFormat)
	downloadPlan, err := preparer.Prepare(ctx, rawURL, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build download arguments: %w", err)
	}

	title := downloadPlan.ExtractInfo.Title

	// If no title, try to get title manually
	if title == "" {
		title, err = d.executor.FetchTitle(
			ctx,
			rawURL,
			idto.WithUseCookies(options.CookieFilePathIfNeeded()),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get title: %w", err)
		}
	}

	// Generate a unique filename if none provided
	uniqueFilePart := options.UniqueFileName
	if uniqueFilePart == "" {
		uniqueFilePart = uuid.New().String()
	}

	var fileName string

	if options.IncludeTitleInFilename {
		fileName = filename.BuildUnique(uniqueFilePart, title, downloadPlan.FileExtNoDot)
	} else {
		fileName = fmt.Sprintf("%s.%s", uniqueFilePart, downloadPlan.FileExtNoDot)
	}

	var (
		fileSize *int64
		tmpSize  int64
	)
	for _, f := range downloadPlan.ExtractInfo.Formats {
		if f.Filesize != nil {
			tmpSize += *f.Filesize
		} else if f.FilesizeApprox != nil {
			tmpSize += *f.FilesizeApprox
		}
	}
	if tmpSize != 0 {
		fileSize = &tmpSize
	}

	channel := d.buildChannel(rawURL, downloadPlan.ExtractInfo)

	meta := &idto.DownloadMeta{
		URL:           rawURL,
		Title:         title,
		Description:   downloadPlan.ExtractInfo.Description,
		FileNameNoExt: filex.FileNameWithoutExt(fileName),
		FileExtNoDot:  downloadPlan.FileExtNoDot,
		FileName:      fileName,
		FileSize:      fileSize,

		Channel: channel,

		MediaInfo: downloadPlan.MediaInfo,
	}

	execOptions := &idto.DownloadExecOptions{
		ConcurrentFragments: options.ConcurrentFragments,
		CookieFilePath:      options.CookieFilePathIfNeeded(),
		Extractor:           downloadPlan.ExtractInfo.Extractor,
		ExtractorArgs:       nil,
		Args:                downloadPlan.Args,
	}

	return meta, execOptions, nil
}
