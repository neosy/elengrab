package httpsrv

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"time"

	authmiddleware "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	"github.com/neosy/elengrab/internal/app/usecases"
	appenv "github.com/neosy/elengrab/internal/pkg/nconfig/app_env"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
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

	// auth
	authMiddleware *authmiddleware.AuthMiddleware

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
		logger:         logger,
		appEnv:         appEnv,
		authMiddleware: authmiddleware.NewAuthMiddleware(logger, deps.Usecases.Auth),
		usecases:       deps.Usecases,
		templates:      deps.Templates,
		assetsDir:      deps.AssetsDir,
		downloadsDir:   deps.DownloadsDir,
	}
}

func (s *httpServer) ListenAndServe(ctx context.Context, port string) error {
	addr := fmt.Sprintf(":%s", port)

	router := s.newRouter()
	handler := nfasthttp.NewHandler(ctx, s.logger, s.appEnv, router.Handler)

	server := &fasthttp.Server{
		Handler: handler,

		// --- Timeouts ---
		// Protects against slowloris during request read
		ReadTimeout: 30 * time.Second,
		// Disable the limit on the duration of recording the response
		// otherwise the long download (10+ minutes) will be interrupted
		WriteTimeout: 0,
		// Limits keep-alive connection lifetime
		IdleTimeout: 60 * time.Second,

		// --- Concurrency and limits ---
		// Maximum number of concurrent connections
		Concurrency: 1024,
		// Prevents a single IP from exhausting workers
		MaxConnsPerIP: 100,
		// Periodically rotates long-lived connections
		MaxRequestsPerConn: 1000,

		// --- Buffers and request sizes ---
		// Sufficient for large headers and cookies
		ReadBufferSize:  32 * 1024,
		WriteBufferSize: 32 * 1024,
		// 10 MB request body limit
		MaxRequestBodySize: 10 * 1024 * 1024,

		// --- Worker pool management ---
		// Cleans up unused workers
		MaxIdleWorkerDuration: 30 * time.Second,

		// --- TCP keep-alive ---
		TCPKeepalive:       true,
		TCPKeepalivePeriod: 30 * time.Second,

		// --- Logging ---
		// Avoid log spam from common network errors
		LogAllErrors: false,
		// Do not log sensitive request data
		SecureErrorLogMessage: true,
	}

	log.Printf("HTTP server listening on %s\n", addr)
	log.Printf("Server started. Run locally: http://localhost:%s (Ctrl+C to stop)", port)

	go func() {
		<-ctx.Done()
		_ = server.Shutdown()
	}()

	err := server.ListenAndServe(addr)
	if err != nil && ctx.Err() == nil {
		s.logger.ErrorContext(ctx, fmt.Sprintf("error run server: %v", err))
		return err
	}

	log.Println("Server stopped")

	return nil
}
