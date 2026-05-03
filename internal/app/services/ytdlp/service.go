package ytdlpsrv

import (
	"log/slog"
	"path/filepath"

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
		if options.YoutubeAllowCookies {
			options.YoutubeAllowCookies = false
			logger.Warn(
				"Deno executable not found in PATH",
				"executable", consts.DenoName,
				"error", err,
			)
			logger.Info("YoutubeAllowCookies has been disabled")
		}
	} else {
		logger.Debug("Deno executable found in PATH", "executable", consts.DenoName)
	}

	if options.YoutubeAllowCookies && options.CookiesDir == "" {
		options.YoutubeAllowCookies = false
		logger.Warn("YoutubeAllowCookies enabled but CookiesDir is empty")
		logger.Info("YoutubeAllowCookies has been disabled")
	}

	if options.YoutubeAllowCookies {
		exists, err := nfile.DirExists(options.CookiesDir)
		if err != nil {
			options.YoutubeAllowCookies = false
			logger.Warn("Error checking if directory exists", "dir", options.CookiesDir, "error", err)
			logger.Info("YoutubeAllowCookies has been disabled")
		} else if !exists {
			options.YoutubeAllowCookies = false
			logger.Warn("Directory cookies does not exist", "dir", options.CookiesDir)
			logger.Info("YoutubeAllowCookies has been disabled")
		} else {
			logger.Debug("Cookies directory exists", "dir", options.CookiesDir)
		}
	}

	if options.YoutubeAllowCookies {
		path := filepath.Join(options.CookiesDir, consts.YtDlpYouTubeCookieFileName)
		exists, err := nfile.FileExists(path)
		if err != nil {
			options.YoutubeAllowCookies = false
			logger.Warn("Failed to check cookies file existence", "path", path, "error", err)
			logger.Info("YoutubeAllowCookies has been disabled")
		} else if !exists {
			options.YoutubeAllowCookies = false
			logger.Warn("Cookies file does not exist", "path", path)
			logger.Info("YoutubeAllowCookies has been disabled")
		} else {
			logger.Debug("Cookies file exists", "path", path)
		}
	}

	if options.YoutubeAllowCookies {
		logger.Info("YoutubeAllowCookies option is enabled")
	}

	return &YtDlpService{
		logger: logger,

		// options
		options: options,

		// internal
		downloader: downloader.NewDownloader(logger, cmdPath, storage, ffmpeg, options),
	}, nil
}
