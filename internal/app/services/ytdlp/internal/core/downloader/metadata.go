package downloader

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/downloader/helper"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/dto"
	"github.com/neosy/elengrab/pkg/nfasthttp"
)

func (d *Downloader) prepareMetadata(
	ctx context.Context,
	url, dlDir, fileName string,
	includeTitleInFilename bool,
	dlOptions dto.DLOptions,
	getBestFormat func(ctx context.Context, url string, format string, useCookies bool) (*idto.MediaInfo, error),
	getTitle func(ctx context.Context, url string, useCookies bool) (string, error),
) (*idto.DownloadMeta, error) {
	// Build yt-dlp arguments and get file extension and title
	args, fileExt, dtoMediaInfo, mediaInfo, err := helper.PrepareDownload(
		ctx,
		url,
		dlOptions,
		getBestFormat,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build download arguments: %w", err)
	}

	// If title is empty, fetch it manually
	title := dtoMediaInfo.Title
	if title == "" {
		var err error
		title, err = getTitle(ctx, url, dlOptions.RequiresYouTubeCookies)
		if err != nil {
			return nil, fmt.Errorf("failed to get title: %w", err)
		}
	}

	// Generate a unique filename if none provided
	if fileName == "" {
		fileName = uuid.New().String()
	}

	if includeTitleInFilename {
		// Sanitize the title to ensure it is a valid file name.
		title := nfasthttp.SanitizeFileName(title)
		// Truncate the title to fit within the maximum length allowed for filenames.
		title = helper.TruncateUTF8(title, consts.MaxTitleLengthInFilename)
		fileName = fmt.Sprintf("%s_%s", nfasthttp.SanitizeFileName(title), fileName)
	}

	fileFullName := fmt.Sprintf("%s.%s", fileName, fileExt)
	filePath := filepath.Join(dlDir, fileFullName)

	var (
		fileSize *int64
		tmpSize  int64
	)
	for _, f := range dtoMediaInfo.Formats {
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
	if dtoMediaInfo.ChannelID != "" {
		channelID = &dtoMediaInfo.ChannelID
	}

	options := idto.DownloadOptions{
		ConcurrentFragments:    dlOptions.ConcurrentFragments,
		RequiresYouTubeCookies: dlOptions.RequiresYouTubeCookies,
		Extractor:              dtoMediaInfo.Extractor,
		ExtractorArgs:          nil,
		Args:                   args,
	}

	return &idto.DownloadMeta{
		URL:          url,
		Title:        title,
		FileName:     fileName,
		FileExt:      fileExt,
		FileFullName: fileFullName,
		FilePath:     filePath,
		FileSize:     fileSize,
		ChannelID:    channelID,
		ChannelURL:   dtoMediaInfo.ChannelUrl,
		ChannelTitle: dtoMediaInfo.ChannelTitle,
		MediaInfo:    mediaInfo,
		Options:      options,
	}, nil
}
