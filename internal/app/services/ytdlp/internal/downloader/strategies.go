package downloader

import (
	"context"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

func (d *Downloader) downloadWithStrategies(
	ctx context.Context,
	url string,
	meta *idto.SafeDownloadMeta,
	onProgressUpdate func(dservices.DownloaderProgress),
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
		attempts = append(attempts,
			downloadAttempt{
				concurrentFragments: new(uint8(1)),
				description:         "single concurrentFragments",
			},
			downloadAttempt{
				concurrentFragments: new(uint8(1)),
				extractorArgs:       new("youtube:player_client=android"),
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
			onProgressUpdate,
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
		isForbidden = isForbidden || helper.HasErrTextForbidden(err)
		if !isForbidden {
			break
		}

	}
	return nil, err
}
