package ytdlpsrv

import (
	"log/slog"

	ffmpegsrv "github.com/neosy/elengrab/internal/app/services/ffmpeg"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/utils"
	nfile "github.com/neosy/elengrab/internal/pkg/file"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

// YtDlpService represents a service for interacting with yt-dlp.
type YtDlpService struct {
	logger *slog.Logger

	// options
	options dto.Options

	// internal
	downloader *downloader.Downloader
}

// NewYtDlpService creates a new YtDlpService instance.
func NewYtDlpService(
	logger *slog.Logger,
	binDir string,
	storage pstorage.DownloadsStorage,
	ffmpeg *ffmpegsrv.FFmpegService,
	opts ...dto.Option,
) (*YtDlpService, error) {
	cmdPath, err := utils.ResolveCmdPath(consts.YtDlpName, binDir)
	if err != nil {
		return nil, err
	}

	options := dto.NewOptions()

	for _, opt := range opts {
		opt(&options)
	}

	_, err = utils.LookupExecutable(consts.DenoName)
	if err != nil {
		if options.AllowCookies {
			options.AllowCookies = false
			logger.Warn(
				"Deno executable not found in PATH",
				"executable", consts.DenoName,
				"error", err,
			)
			logger.Info("AllowCookies has been disabled")
		}
	} else {
		logger.Debug("Deno executable found in PATH", "executable", consts.DenoName)
	}

	if options.AllowCookies && options.CookiesDir == "" {
		options.AllowCookies = false
		logger.Warn("AllowCookies enabled but CookiesDir is empty")
		logger.Info("AllowCookies has been disabled")
	}

	if options.AllowCookies {
		exists, err := nfile.DirExists(options.CookiesDir)
		if err != nil {
			options.AllowCookies = false
			logger.Warn("Error checking if directory exists", "dir", options.CookiesDir, "error", err)
			logger.Info("AllowCookies has been disabled")
		} else if !exists {
			options.AllowCookies = false
			logger.Warn("Directory cookies does not exist", "dir", options.CookiesDir)
			logger.Info("AllowCookies has been disabled")
		} else {
			logger.Debug("Cookies directory exists", "dir", options.CookiesDir)
		}
	}

	if options.AllowCookies {
		logger.Info("AllowCookies option is enabled")
	}

	return &YtDlpService{
		logger: logger,

		// options
		options: options,

		// internal
		downloader: downloader.NewDownloader(logger, cmdPath, storage, ffmpeg, options),
	}, nil
}
