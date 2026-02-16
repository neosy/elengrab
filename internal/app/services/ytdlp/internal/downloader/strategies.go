package downloader

import (
	"context"
	"strings"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (d *Downloader) downloadWithStrategies(
	ctx context.Context,
	url string,
	meta *idto.SafeDownloadMeta,
	sendData func(*ddownload.DownloadResult),
) ([]byte, error) {
	// Define download attempts
	type downloadAttempt struct {
		extractorArgs       *string
		concurrentFragments *uint8
		description         string
	}

	var attempts []downloadAttempt

	// default
	attempts = append(attempts, downloadAttempt{
		description: "default",
	})

	// Additional attempts for YouTube extractor
	if meta.Meta.Options.Extractor == "youtube" {
		one := uint8(1)
		android := "youtube:player_client=android"

		attempts = append(attempts,
			downloadAttempt{
				concurrentFragments: &one,
				description:         "single concurrentFragments",
			},
			downloadAttempt{
				concurrentFragments: &one,
				extractorArgs:       &android,
				description:         "android client",
			},
		)
	}

	// Helper function to wait for a specified duration or context cancellation
	waitOrCancel := func(delay time.Duration) error {
		d.logger.Debug("Waiting...", "duration", delay.String())
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			return nil
		}
	}

	var (
		out         []byte
		err         error
		isForbidden bool
	)

	// Try download attempts
	for i, attempt := range attempts {
		if i > 0 {
			if err := waitOrCancel(consts.YtDlpRetryDelay); err != nil {
				return nil, err
			}
		}

		d.logger.Debug(
			"Download attempt",
			"strategy", attempt.description,
			"url", url,
		)

		// Create a copy of the metadata for this attempt
		metaCopy := meta.CopyMeta()
		if attempt.extractorArgs != nil {
			metaCopy.Options.ExtractorArgs = attempt.extractorArgs
		}
		if attempt.concurrentFragments != nil {
			metaCopy.Options.ConcurrentFragments = *attempt.concurrentFragments
		}

		// Function to run the download process
		// Run yt-dlp for the given URL and metadata.
		// Capture output and error; send error if execution fails.
		out, err = d.executor.RunYtDlp(
			ctx,
			url,
			metaCopy,
			func(progress ddownload.DownloadProgress) {
				meta.Lock()
				meta.Meta.Progress = &progress
				meta.Unlock()
				sendData(meta.InitialResult())
			},
		)
		if err == nil {
			return out, nil
		}

		d.logger.Warn(
			"Download failed",
			"url", url,
			"strategy", attempt.description,
			"error", err,
		)

		// If not a 403 error, do not retry
		isForbidden = isForbidden || d.isForbidden(err)
		if !isForbidden {
			break
		}

	}
	return nil, err
}

func (d *Downloader) isForbidden(err error) bool {
	return err != nil && strings.Contains(err.Error(), "HTTP Error 403")
}
