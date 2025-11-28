package httpsrv

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases"
	"github.com/valyala/fasthttp"
)

type Dependencies struct {
	Usecases  *usecases.Usecases
	Templates *template.Template

	// Options
	AssetsDir string
}

type httpServer struct {
	logger   *slog.Logger
	usecases *usecases.Usecases

	templates *template.Template

	// Options
	assetsDir string
}

func NewServer(logger *slog.Logger, deps *Dependencies) *httpServer {
	return &httpServer{
		logger:    logger,
		usecases:  deps.Usecases,
		templates: deps.Templates,
		assetsDir: deps.AssetsDir,
	}
}

func (s *httpServer) ListenAndServe(ctx context.Context, port string) error {
	addr := fmt.Sprintf(":%s", port)

	router := s.newRouter()
	handler := router.Handler

	log.Printf("HTTP server listening on %s\n", addr)

	err := fasthttp.ListenAndServe(addr, handler)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error run server: %v", err))
		return err
	}

	return nil
}
