package httpsrv

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/middleware"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/routes"
	"github.com/neosy/elengrab/internal/app/usecases"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/infrastructure/observability/metrics"
	appenv "github.com/neosy/elengrab/internal/pkg/config/app_env"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
	"github.com/valyala/fasthttp"
)

type Dependencies struct {
	// Storages
	DownloadsStorage pstorage.DownloadsStorage

	// caches
	AssetFileCacheRepository persistence.AssetFileCacheRepository

	Usecases *usecases.Usecases
	Template *template.Template

	// Options
	AppMode         dtypes.AppMode
	BaseURL         string
	ShortLinkPrefix string
	AssetsDir       string
	MetricsEnabled  bool
}

type httpServer struct {
	logger *slog.Logger
	appEnv appenv.AppEnv

	// middleware
	middlewares *middleware.Middlewares

	handlers *handlers.Handlers

	// usecases
	usecases *usecases.Usecases

	// templates
	template *template.Template

	// Options
	appMode         dtypes.AppMode
	baseURL         string
	shortLinkPrefix string
	assetsDir       string
	metricsEnabled  bool
}

func NewServer(logger *slog.Logger, appEnv appenv.AppEnv, deps *Dependencies) (*httpServer, error) {
	handlers, err := handlers.New(
		logger,
		&handlers.Dependencies{
			DownloadsStorage: deps.DownloadsStorage,
			Assets:           assets.NewAssets(deps.AssetsDir, deps.AssetFileCacheRepository),
			Usecases:         deps.Usecases,
			Template:         deps.Template,
			AppEnv:           appEnv,
			AppMode:          deps.AppMode,
			BaseURL:          deps.BaseURL,
			ShortLinkPrefix:  deps.ShortLinkPrefix,
		},
	)
	if err != nil {
		return nil, err
	}

	return &httpServer{
		logger: logger,
		appEnv: appEnv,

		middlewares: middleware.NewMiddlewares(
			logger,
			deps.Usecases,
			handlers,
			deps.AppMode,
		),

		handlers: handlers,

		usecases: deps.Usecases,
		template: deps.Template,

		// Optons
		appMode:         deps.AppMode,
		baseURL:         deps.BaseURL,
		shortLinkPrefix: deps.ShortLinkPrefix,
		assetsDir:       deps.AssetsDir,
		metricsEnabled:  deps.MetricsEnabled,
	}, nil
}

func (s *httpServer) ListenAndServe(ctx context.Context, port string) error {
	addr := fmt.Sprintf(":%s", port)

	routes := routes.NewRoutes(s.middlewares, s.handlers, s.shortLinkPrefix)
	router := routes.Router()
	handler := nfasthttp.NewHandler(ctx, s.logger, s.appEnv, router.Handler)

	if s.metricsEnabled {
		handler = metrics.MiddlewareHandler(handler)
	}

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

		DisableKeepalive: false,

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
