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
	dlOptions :=
		helper.PrepareDownloadOptions(url, d.serviceOptions, options)

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

	var needCookies bool
	err = d.executor.EnsureFormatCache(ctx, url)
	if err != nil {
		needCookies = dlOptions.AllowCookies() && helper.CheckCookiesError(err)
		if !needCookies {
			sendError(nil, errorx.Errorf("failed to ensure format cache: %w", err))
			return
		}
	}

	if needCookies {
		d.logger.Debug(
			"Cookies are enabled",
			"title", title,
			"url", url,
		)

		err = d.executor.EnsureFormatCache(ctx, url, idto.WithUseCookies(dlOptions.CookieFilePath))
		if err != nil {
			sendError(nil, errorx.Errorf("failed to ensure format cache: %w", err))
			return
		}

		dlOptions.NeedsCookies = true
	}

	// Prepare download metadata
	var (
		meta        idto.SafeDownloadMeta
		execOptions *idto.DownloadExecOptions
	)
	meta.Meta, execOptions, err = d.prepareDownload(
		ctx,
		url,
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
			channel := d.fetchAndBuildChannelAvatar(meta.CopyMeta())
			if channel != nil {
				meta.Lock()
				meta.Meta.Channel = channel
				meta.Unlock()
				sendData(meta.InitialResult())
			}
		})
	}

	// Start asynchronous get of thumbnail
	wg.Go(func() {
		ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
		defer cancel()
		imgData := d.extractThumbnailFromURL(
			ctx, url,
			false, // Not needed, loading from cache.
		)
		if imgData != nil {
			meta.Lock()
			meta.Meta.Thumbnail = imgData
			meta.Unlock()
			sendData(meta.InitialResult())
		}
	})

	out, err := d.downloadWithStrategies(
		ctx, url, meta.CopyMeta(), execOptions.Copy(),
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

	// Get (asynchronous) the actual file size
	wg.Go(func() {
		size := d.fileSize(meta.Meta.FileFullName)
		if size != 0 {
			meta.Lock()
			meta.Meta.FileSize = &size
			meta.Unlock()
		}
	})

	// Start asynchronous reading metadata from the given file.
	wg.Go(func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		mediaInfo := d.extractMediaInfoFromFile(ctx, meta.Meta.FileFullName, meta.Meta.MediaInfo)
		if mediaInfo != nil {
			meta.Lock()
			if mediaInfo.Duration > 0 {
				meta.Meta.MediaInfo.Duration = mediaInfo.Duration
			}
			if mediaInfo.VideoInfo != nil {
				meta.Meta.MediaInfo.VideoInfo = mediaInfo.VideoInfo
			}
			if mediaInfo.AudioInfo != nil {
				meta.Meta.MediaInfo.AudioInfo = mediaInfo.AudioInfo
			}
			meta.Unlock()
			sendData(meta.InitialResult())
		}
	})

	// Start asynchronous extracting a thumbnail frame from the video file.
	if meta.Meta.MediaInfo.HasVideo() {
		wg.Go(func() {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			imgData := d.extractThumbnailFromFile(ctx, meta.Meta.FileFullName)
			if imgData != nil {
				meta.Lock()
				meta.Meta.ThumbnailVideoFrame = imgData
				meta.Unlock()
				sendData(meta.InitialResult())
			}
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

func (d *Downloader) extractThumbnailFromFile(ctx context.Context, fileName string) *dtypes.ImageData {
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
	imgData, err := extractBestFrame(ctx, d.storage.Path(fileName))
	if err != nil {
		return nil
	}

	return imgData
}

func (d *Downloader) extractMediaInfoFromFile(
	ctx context.Context, fileName string, mediaInfo *dservices.MediaInfo,
) *dservices.MediaInfo {
	mediaInfo, err := d.ffmpeg.ExtractVideoAudioInfoFromFile(ctx, d.storage.Path(fileName), mediaInfo)
	if err != nil {
		d.logger.Warn(
			"Failed to get media info from file",
			"filePath", d.storage.Path(fileName),
			"error", err,
		)
		return nil
	}

	return mediaInfo
}

func (d *Downloader) partialHash(filePath string) *string {
	h, err := hash.FilePartialHash(filePath)
	if err == nil && h != "" {
		return &h
	}
	return nil
}

func (d *Downloader) fileSize(fileName string) int64 {
	fileInfo, err := os.Stat(d.storage.Path(fileName))
	if err != nil {
		return 0
	}
	return int64(fileInfo.Size())
}
