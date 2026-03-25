package ytdlpsrv

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/ffmpeg"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/utils"
	"github.com/neosy/elengrab/internal/pkg/nfile"
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
	downloadsDir string,
	options *dto.Options,
) (*YtDlpService, error) {
	cmdPath, err := utils.ResolveCmdPath(consts.YtDlpName, binDir)
	if err != nil {
		return nil, err
	}

	err = ffmpeg.CheckFFmpeg(consts.FFmpegName)
	if err != nil {
		return nil, err
	} else {
		logger.Debug("FFmpeg executable found in PATH", "executable", consts.FFmpegName)
	}

	_, err = utils.LookupExecutable(consts.DenoName)
	if err != nil {
		if options != nil && options.YoutubeAllowCookies {
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

	if downloadsDir == "" {
		return nil, errors.New("download directory is not set")
	}

	if strings.HasSuffix(downloadsDir, "/") || strings.HasSuffix(downloadsDir, "\\") {
		return nil, fmt.Errorf("downloads directory must not end with a slash or backslash: %s", downloadsDir)
	}

	if err := nfile.CheckDir(downloadsDir); err != nil {
		return nil, err
	}

	var opts dto.Options
	if options != nil {
		opts = *options
	}
	opts.SetDefaults(false)

	if opts.YoutubeAllowCookies && opts.CookiesDir == "" {
		opts.YoutubeAllowCookies = false
		logger.Warn("YoutubeAllowCookies enabled but CookiesDir is empty")
		logger.Info("YoutubeAllowCookies has been disabled")
	}

	if opts.YoutubeAllowCookies {
		exists, err := nfile.DirExists(opts.CookiesDir)
		if err != nil {
			opts.YoutubeAllowCookies = false
			logger.Warn("Error checking if directory exists", "dir", opts.CookiesDir, "error", err)
			logger.Info("YoutubeAllowCookies has been disabled")
		} else if !exists {
			opts.YoutubeAllowCookies = false
			logger.Warn("Directory cookies does not exist", "dir", opts.CookiesDir)
			logger.Info("YoutubeAllowCookies has been disabled")
		} else {
			logger.Debug("Cookies directory exists", "dir", opts.CookiesDir)
		}
	}

	if opts.YoutubeAllowCookies {
		path := filepath.Join(opts.CookiesDir, consts.YtDlpYouTubeCookieFileName)
		exists, err := nfile.FileExists(path)
		if err != nil {
			opts.YoutubeAllowCookies = false
			logger.Warn("Failed to check cookies file existence", "path", path, "error", err)
			logger.Info("YoutubeAllowCookies has been disabled")
		} else if !exists {
			opts.YoutubeAllowCookies = false
			logger.Warn("Cookies file does not exist", "path", path)
			logger.Info("YoutubeAllowCookies has been disabled")
		} else {
			logger.Debug("Cookies file exists", "path", path)
		}
	}

	if opts.YoutubeAllowCookies {
		logger.Info("YoutubeAllowCookies option is enabled")
	}

	return &YtDlpService{
		logger: logger,

		// options
		options: opts,

		// internal
		downloader: downloader.NewDownloader(logger, cmdPath, downloadsDir, options),
	}, nil
}
