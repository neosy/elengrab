package downloader

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	downloadpreparer "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/download_preparer"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/stringx"
)

func (d *Downloader) prepareDownload(
	ctx context.Context,
	url string,
	options idto.DLOptions,
) (*idto.DownloadMeta, *idto.DownloadExecOptions, error) {
	// Build yt-dlp arguments and get file extension and title
	preparer := downloadpreparer.NewDownloadPreparer(d.executor.FetchInfoWithBestFormat)
	downloadPlan, err := preparer.Prepare(ctx, url, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build download arguments: %w", err)
	}

	title := downloadPlan.ExtractInfo.Title

	// If no title, try to get title manually
	if title == "" {
		title, err = d.executor.FetchTitle(
			ctx,
			url,
			idto.WithUseCookies(options.CookieFilePathIfNeeded()),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get title: %w", err)
		}
	}

	// Generate a unique filename if none provided
	fileName := options.FileName
	if fileName == "" {
		fileName = uuid.New().String()
	}

	if options.IncludeTitleInFilename {
		// Sanitize the title to ensure it is a valid file name.
		title := nfasthttp.SanitizeFileName(title)
		// Truncate the title to fit within the maximum length allowed for filenames.
		title = stringx.TruncateBytesWords(title, consts.MaxTitleLengthInFilename)

		fileName = fmt.Sprintf("%s_%s", title, fileName)
	}

	fileFullName := fmt.Sprintf("%s.%s", fileName, downloadPlan.FileExt)

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

	var channelID *string
	if downloadPlan.ExtractInfo.ChannelID != "" {
		channelID = &downloadPlan.ExtractInfo.ChannelID
	}

	meta := &idto.DownloadMeta{
		URL:          url,
		Title:        title,
		Description:  downloadPlan.ExtractInfo.Description,
		FileName:     fileName,
		FileExt:      downloadPlan.FileExt,
		FileFullName: fileFullName,
		FileSize:     fileSize,
		ChannelID:    channelID,
		ChannelURL:   downloadPlan.ExtractInfo.ChannelUrl,
		ChannelTitle: downloadPlan.ExtractInfo.ChannelTitle,
		MediaInfo:    downloadPlan.MediaInfo,
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
