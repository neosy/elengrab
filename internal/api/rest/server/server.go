package httpsrv

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases"
	appenv "github.com/neosy/elengrab/pkg/nconfig/app_env"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

type Dependencies struct {
	Usecases  *usecases.Usecases
	Templates *template.Template

	// Options
	AssetsDir    string
	DownloadsDir string
}

type httpServer struct {
	logger *slog.Logger
	appEnv appenv.AppEnv

	// usecases
	usecases *usecases.Usecases

	// templates
	templates *template.Template

	// Options
	assetsDir    string
	downloadsDir string
}

func NewServer(logger *slog.Logger, appEnv appenv.AppEnv, deps *Dependencies) *httpServer {
	return &httpServer{
		logger:       logger,
		appEnv:       appEnv,
		usecases:     deps.Usecases,
		templates:    deps.Templates,
		assetsDir:    deps.AssetsDir,
		downloadsDir: deps.DownloadsDir,
	}
}

func (s *httpServer) ListenAndServe(ctx context.Context, port string) error {
	addr := fmt.Sprintf(":%s", port)

	router := s.newRouter()
	handler := nfasthttp.NewHandler(ctx, s.logger, s.appEnv, router.Handler)

	log.Printf("HTTP server listening on %s\n", addr)

	err := fasthttp.ListenAndServe(addr, handler)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Sprintf("error run server: %v", err))
		return err
	}

	return nil
}
