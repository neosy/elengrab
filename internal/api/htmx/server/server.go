package httpsrv

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/blocktree/openwallet/v2/log"
	"github.com/neosy/elengrab/internal/app/usecases"
	"github.com/valyala/fasthttp"
)

type Dependencies struct {
	Usecases *usecases.Usecases

	// Options
	AssetsDir    string
	DownloadsDir string
}

type httpServer struct {
	usecases *usecases.Usecases

	// Options
	assetsDir    string
}

func NewServer(logger *slog.Logger, deps *Dependencies) *httpServer {
	return &httpServer{
		usecases:     deps.Usecases,
		assetsDir:    deps.AssetsDir,
	}
}

func (s *httpServer) ListenAndServe(ctx context.Context, port string) error {
	addr := fmt.Sprintf(":%s", port)

	router := s.newRouter(s.assetsDir)
	handler := router.Handler

	log.Info(fmt.Sprintf("HTTP server listening on %s", addr))

	err := fasthttp.ListenAndServe(addr, handler)
	if err != nil {
		log.Error(fmt.Sprintf("error run server: %v", err))
		return err
	}

	return nil
}
