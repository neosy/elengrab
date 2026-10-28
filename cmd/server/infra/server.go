package infra

import (
	"log/slog"

	httpsrv "github.com/neosy/elengrab/internal/api/rest/server"
	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/app"
	upaths "github.com/neosy/elengrab/internal/app/utils/paths"
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func StartHTTPServer(logger *slog.Logger, cfg *iconfig.Config, app *app.Application, ready chan<- struct{}) {
	tmpl, err := assets.LoadTemplates(upaths.AbsPath(cfg.Elengrab.RootDir, cfg.Elengrab.AssetsDir))
	if err != nil {
		logger.Error("Failed to load templates", "error", err)
		app.Cancel()
		return
	}

	deps := &httpsrv.Dependencies{
		DownloadsStorage:      app.DownloadsStorage,
		AssetFileCacheRepo:    app.AssetFileCacheRepo,
		UploadedFileCacheRepo: app.UploadedFileCacheRepo,
		SSEEventLogCacheRepo:  app.SSEEventLogCacheRepo,
		Usecases:              app.Usecases,
		Template:              tmpl,

		// Options
		AppMode:         dtypes.MustParseAppMode(cfg.Elengrab.Mode),
		BaseURL:         cfg.Elengrab.BaseURL,
		ShortLinkPrefix: cfg.Elengrab.ShortLinkPrefix,
		AssetsDir:       upaths.AbsPath(cfg.Elengrab.RootDir, cfg.Elengrab.AssetsDir),
		MetricsEnabled:  cfg.AdminServer.Enable && cfg.AdminServer.DebugConfig.EnableMetrics,
		MaxUploadSize:   cfg.Elengrab.Upload.MaxFileSize(),
	}

	httpServer, err := httpsrv.NewServer(logger, cfg.AppConfig.AppEnv, deps)
	if err != nil {
		logger.Error("Failed to create HTTP server", "error", err)
		app.Cancel()
		return
	}

	err = httpServer.ListenAndServe(app.Context(), cfg.HTTPServer.Port, ready)
	if err != nil {
		app.Cancel()
	}
}
