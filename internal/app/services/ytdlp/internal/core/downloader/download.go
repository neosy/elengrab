package downloader

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/downloader/ffmpeg"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/downloader/helper"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/dto"
	"github.com/neosy/elengrab/internal/app/utils"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	"github.com/neosy/elengrab/pkg/nfile"
	"github.com/neosy/elengrab/pkg/syncx"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (d *Downloader) Download(
	ctx context.Context,
	url string,
	concurrentFragments uint8,
	options *dservices.DownloadOptions,
	getBestFormat func(ctx context.Context, url string, format string) (*idto.MediaInfo, error),
	getTitle func(ctx context.Context, url string) (string, error),
	downloadResultCh chan<- *ddownload.DownloadResult,
) {
	var wg sync.WaitGroup
	done := syncx.NewDoneSignal()
	defer func() {
		done.Close()
		wg.Wait()
	}()

	var (
		sendError = func(data *ddownload.DownloadResult, err error) {
			var result = &ddownload.DownloadResult{}
			if data != nil {
				*result = *data
			}
			result.Error = err

			select {
			case <-done.Done():
			case downloadResultCh <- result:
			case <-ctx.Done():
			}
		}

		sendData = func(data *ddownload.DownloadResult) {
			select {
			case <-done.Done():
			case downloadResultCh <- data:
			case <-ctx.Done():
			}
		}
	)

	url = strings.TrimSpace(url)

	// Prepare download options with defaults and user overrides
	dlOptions, dlDir, fileName, includeTitleInFilename :=
		helper.PrepareDownloadOptions(d.downloadsDir, concurrentFragments, options)

	// Ensure download directory exists
	if err := nfile.CheckDir(dlDir); err != nil {
		sendError(nil, fmt.Errorf("failed to check directory: %w", err))
		return
	}

	// Try to get the title fast
	title, err := helper.GetTitleFast(url)
	if err != nil {
		d.logger.Debug(
			"Get title fast",
			"url", url,
			"error", err,
		)
	}
	if err == nil && title != "" {
		d.logger.Debug(
			"Get title fast",
			"url", url,
			"title", title,
		)
		sendData(&ddownload.DownloadResult{
			MediaTitle: title,
		})
	}

	// Prepare download metadata
	var meta idto.SafeDownloadMeta
	meta.Meta, err = d.prepareMetadata(
		ctx,
		url,
		dlDir,
		fileName,
		includeTitleInFilename,
		dlOptions,
		getBestFormat,
		getTitle,
	)
	if err != nil {
		sendError(&ddownload.DownloadResult{MediaTitle: title}, err)
		return
	}

	sendData(meta.InitialResult())

	// Start asynchronous fetching of the channel avatar.
	// Returns a channel from which the avatar can be read once the goroutine completes.
	d.fetchChannelAvatarAsync(
		&wg,
		meta.CopyMeta(),
		func(avatar *ddownload.DownloadResultChannelAvatar) {
			meta.Lock()
			meta.Meta.ChannelAvatar = avatar
			meta.Unlock()
			sendData(meta.InitialResult())
		},
	)

	out, err := d.downloadWithStrategies(
		ctx,
		url,
		&meta,
		sendData,
	)
	if err != nil {
		// Deleting the cache, because Youtube could have changed the format
		d.formatCache.DeleteByURL(url)
		// Send error result
		sendError(meta.InitialResult(), err)
		return
	}

	// Debug log command output
	d.logger.Debug(
		"Download completed",
		"url", url,
		"out", string(out),
	)

	// Get the actual file size
	fileInfo, err := os.Stat(meta.Meta.FilePath)
	if err == nil {
		meta.Lock()
		meta.Meta.FileSize = uptr.Int64(fileInfo.Size())
		meta.Unlock()
	}

	var partialHash *string
	{
		h, err := utils.HashPartialMedia(meta.Meta.FilePath)
		if err == nil && h != "" {
			partialHash = &h
		}
	}

	// Waiting for background processes to complete
	wg.Wait()

	// Reading metadata from the given file.
	vInfo, aInfo := ffmpeg.NewInfo().GetVideoAudioInfoFromFile(ctx, meta.Meta.FilePath, meta.Meta.MediaInfo)
	if vInfo != nil {
		meta.Lock()
		meta.Meta.MediaInfo.VideoInfo = vInfo
		meta.Unlock()
	}
	if aInfo != nil {
		meta.Lock()
		meta.Meta.MediaInfo.AudioInfo = aInfo
		meta.Unlock()
	}

	// Build response struct
	result := meta.InitialResult()
	result.PartialHash = partialHash

	// Build response struct
	sendData(result)

	// Log successful download
	// Info Download completed
	d.logger.Info(
		"Download completed",
		"title", meta.Meta.Title,
		"url", url,
		"mediaInfo", result.MediaInfo,
	)
	d.logger.Debug("Download success", "meta", meta.Meta)
}
