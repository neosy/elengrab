package downloader

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/ffmpeg"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	"github.com/neosy/elengrab/internal/app/utils"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	"github.com/neosy/elengrab/pkg/nfile"
	"github.com/neosy/elengrab/pkg/syncx"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (d *Downloader) Download(
	ctx context.Context,
	url string,
	options *dservices.DownloadOptions,
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
		helper.PrepareDownloadOptions(d.downloadsDir, d.serviceOptions.ConcurrentFragments, options)

	// Ensure download directory exists
	if err := nfile.CheckDir(dlDir); err != nil {
		sendError(nil, fmt.Errorf("failed to check directory: %w", err))
		return
	}

	// Try to fetch the title fast
	startTime := time.Now()
	title, err := helper.FetchTitleFast(ctx, url)
	elapsed := time.Since(startTime)
	if err != nil {
		d.logger.Debug(
			"Failed to fetch title fast",
			"url", url,
			"error", err,
		)
	}
	if err == nil && title != "" {
		d.logger.Debug(
			"Fetched title (fast)",
			"title", title,
			"url", url,
			"elapsed", uformat.DurationFormat(elapsed),
		)
		sendData(&ddownload.DownloadResult{
			MediaTitle: title,
		})
	}

	var requiresYouTubeCookies = false
	err = d.executor.EnsureFormatCache(ctx, url, false)
	if err != nil {
		requiresYouTubeCookies = d.serviceOptions.YoutubeAllowCookies &&
			helper.CheckForYouTubeCookiesError(err)
		if !requiresYouTubeCookies {
			sendError(nil, fmt.Errorf("failed to ensure format cache: %w", err))
			return
		}
	}

	if requiresYouTubeCookies {
		d.logger.Debug(
			"YouTube cookies are enabled",
			"title", title,
			"url", url,
		)
		err = d.executor.EnsureFormatCache(ctx, url, true)
		if err != nil {
			sendError(nil, fmt.Errorf("failed to ensure format cache: %w", err))
			return
		}
	}

	dlOptions.RequiresYouTubeCookies = requiresYouTubeCookies

	// Prepare download metadata
	var meta idto.SafeDownloadMeta
	meta.Meta, err = d.prepareMetadata(
		ctx,
		url,
		dlDir,
		fileName,
		includeTitleInFilename,
		dlOptions,
	)
	if err != nil {
		sendError(&ddownload.DownloadResult{MediaTitle: title}, err)
		return
	}

	sendData(meta.InitialResult())

	// Start asynchronous fetching of the channel avatar.
	// Returns a channel from which the avatar can be read once the goroutine completes.
	if options.DownloadChannelAvatar {
		d.fetchChannelAvatarAsync(
			&wg,
			meta.CopyMeta(),
			func(channel *ddownload.DownloadChannel) {
				meta.Lock()
				meta.Meta.Channel = channel
				meta.Unlock()
				sendData(meta.InitialResult())
			},
		)
	}

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
