package downloader

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	ffmpegsrv "github.com/neosy/elengrab/internal/app/services/ffmpeg"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	"github.com/neosy/elengrab/internal/app/utils/hash"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/syncx"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (d *Downloader) Download(
	ctx context.Context,
	url string,
	options *dservices.DownloadOptions,
	downloadResultCh chan<- *dservices.DownloaderResult,
) {
	var wg sync.WaitGroup
	done := syncx.NewDoneSignal()
	defer func() {
		done.Close()
		wg.Wait()
	}()

	var (
		sendError = func(data *dservices.DownloaderResult, err error) {
			var result = &dservices.DownloaderResult{}
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

		sendData = func(data *dservices.DownloaderResult) {
			select {
			case <-done.Done():
			case downloadResultCh <- data:
			case <-ctx.Done():
			}
		}
	)

	url = strings.TrimSpace(url)

	// Prepare download options with defaults and user overrides
	dlOptions, fileName, includeTitleInFilename :=
		helper.PrepareDownloadOptions(d.serviceOptions.ConcurrentFragments, options)

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
		sendData(&dservices.DownloaderResult{
			MediaTitle: title,
		})
	}

	var requiresYouTubeCookies = false
	err = d.executor.EnsureFormatCache(ctx, url, false)
	if err != nil {
		requiresYouTubeCookies = d.serviceOptions.YoutubeAllowCookies &&
			helper.CheckForYouTubeCookiesError(err)
		if !requiresYouTubeCookies {
			sendError(nil, errorx.Errorf("failed to ensure format cache: %w", err))
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
			sendError(nil, errorx.Errorf("failed to ensure format cache: %w", err))
			return
		}
	}

	dlOptions.RequiresYouTubeCookies = requiresYouTubeCookies

	// Prepare download metadata
	var meta idto.SafeDownloadMeta
	meta.Meta, err = d.prepareMetadata(
		ctx,
		url,
		fileName,
		includeTitleInFilename,
		dlOptions,
	)
	if err != nil {
		sendError(&dservices.DownloaderResult{MediaTitle: title}, err)
		return
	}

	sendData(meta.InitialResult())

	// Start asynchronous fetching of the channel avatar.
	// Returns a channel from which the avatar can be read once the goroutine completes.
	if options.DownloadChannelAvatar {
		wg.Go(func() {
			d.fetchAndBuildChannelAvatar(
				meta.CopyMeta(),
				func(channel *dtypes.Channel) {
					meta.Lock()
					meta.Meta.Channel = channel
					meta.Unlock()
					sendData(meta.InitialResult())
				},
			)
		})
	}

	// Start asynchronous get of thumbnail
	wg.Go(func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		d.extractThumbnailFromURL(
			ctx, url,
			false, // Not needed, loading from cache.
			func(imageData *dtypes.ImageData) {
				meta.Lock()
				meta.Meta.Thumbnail = imageData
				meta.Unlock()
				sendData(meta.InitialResult())
			},
		)
	})

	out, err := d.downloadWithStrategies(
		ctx, url, &meta,
		func(progress dservices.DownloaderProgress) {
			meta.Lock()
			meta.Meta.Progress = &progress
			meta.Unlock()
			sendData(meta.InitialResult())
		})
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
	d.refreshFileSize(&meta)

	// Waiting for background processes to complete
	wg.Wait()

	// Start asynchronous reading metadata from the given file.
	wg.Go(func() {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		d.refreshMediaInfo(ctx, &meta)
		sendData(meta.InitialResult())
	})

	// Start asynchronous extracting a thumbnail frame from the video file.
	if meta.Meta.MediaInfo.HasVideo() {
		wg.Go(func() {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			d.extractAndSaveThumbnail(ctx, &meta)
			sendData(meta.InitialResult())
		})
	}

	// Waiting for background processes to complete
	wg.Wait()

	// Build response struct
	result := meta.InitialResult()
	result.PartialHash = d.partialHash(d.storage.Path(meta.Meta.FileFullName))

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

func (d *Downloader) extractAndSaveThumbnail(ctx context.Context, meta *idto.SafeDownloadMeta) {
	extractBestFrame := func(ctx context.Context, filePath string) (*dtypes.ImageData, error) {
		imgData, err := d.ffmpeg.ExtractBestFrame(
			ctx,
			filePath,
			ffmpegsrv.WithFrameStrategy(ffmpegsrv.FrameStrategyBalanced{}),
			ffmpegsrv.WithFrameFormat(ffmpegsrv.FrameFormatWebP{}),
		)
		if err != nil {
			imgData, err = d.ffmpeg.ExtractBestFrame(ctx, filePath)
		}

		if err != nil {
			d.logger.Warn(
				"Failed to get thumbnail from file",
				"filePath", filePath,
				"error", err,
			)
			return nil, err
		}

		return imgData, nil
	}

	var err error
	imgData, err := extractBestFrame(ctx, d.storage.Path(meta.Meta.FileFullName))
	if err != nil {
		return
	}
	if imgData != nil {
		meta.Lock()
		meta.Meta.ThumbnailVideoFrame = imgData
		meta.Unlock()
	}
}

func (d *Downloader) refreshMediaInfo(ctx context.Context, meta *idto.SafeDownloadMeta) {
	vInfo, aInfo, err := d.ffmpeg.GetVideoAudioInfoFromFile(ctx, d.storage.Path(meta.Meta.FileFullName), meta.Meta.MediaInfo)
	if err != nil {
		d.logger.Warn(
			"Failed to get media info from file",
			"filePath", d.storage.Path(meta.Meta.FileFullName),
			"error", err,
		)
		return
	}
	if vInfo != nil || aInfo != nil {
		meta.Lock()
		if vInfo != nil {
			meta.Meta.MediaInfo.VideoInfo = vInfo
		}
		if aInfo != nil {
			meta.Meta.MediaInfo.AudioInfo = aInfo
		}
		meta.Unlock()
	}
}

func (d *Downloader) partialHash(filePath string) *string {
	h, err := hash.FilePartialHash(filePath)
	if err == nil && h != "" {
		return &h
	}
	return nil
}

func (d *Downloader) refreshFileSize(meta *idto.SafeDownloadMeta) {
	fileInfo, err := os.Stat(d.storage.Path(meta.Meta.FileFullName))
	if err == nil {
		meta.Lock()
		meta.Meta.FileSize = uptr.Int64(fileInfo.Size())
		meta.Unlock()
	}
}
