package infra

import (
	"log/slog"

	httpsrv "github.com/neosy/elengrab/internal/api/rest/server"
	httptemplates "github.com/neosy/elengrab/internal/api/rest/server/templates"
	"github.com/neosy/elengrab/internal/app"
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func StartHTTPServer(logger *slog.Logger, cfg *iconfig.Config, app *app.Application) {
	tmpl, err := httptemplates.LoadTemplates(absPath(cfg.Elengrab.RootDir, cfg.Elengrab.AssetsDir))
	if err != nil {
		logger.Error("Failed to load templates", "error", err)
		app.Cancel()
		return
	}

	deps := &httpsrv.Dependencies{
		Usecases:        app.Usecases,
		Templates:       tmpl,
		AppMode:         dtypes.MustParseAppMode(cfg.Elengrab.Mode),
		BaseURL:         cfg.Elengrab.BaseURL,
		ShortLinkPrefix: cfg.Elengrab.ShortLinkPrefix,
		AssetsDir:       absPath(cfg.Elengrab.RootDir, cfg.Elengrab.AssetsDir),
		DownloadsDir:    absPath(cfg.Elengrab.RootDir, cfg.Elengrab.DownloadsDir),
	}

	httpServer := httpsrv.NewServer(logger, cfg.AppConfig.AppEnv, deps)
	err = httpServer.ListenAndServe(app.Context(), cfg.HTTPServer.Port)
	if err != nil {
		app.Cancel()
	}
}
