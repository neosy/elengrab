package httpsrv

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/middleware"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
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

const (
	uploadMultipartOverhead       = 1 << 20  // multipart boundaries and part headers
	defaultMaxUploadSize    int64 = 10 << 30 // 10 GiB
)

type Dependencies struct {
	// Storages
	DownloadsStorage pstorage.DownloadsStorage

	// Caches
	AssetFileCacheRepo    persistence.AssetFileCacheRepository
	UploadedFileCacheRepo persistence.UploadedFileCacheRepository

	Usecases *usecases.Usecases
	Template *template.Template

	// Options
	AppMode         dtypes.AppMode
	BaseURL         string
	ShortLinkPrefix string
	AssetsDir       string
	MetricsEnabled  bool
	MaxUploadSize   int64
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
	maxUploadSize   int64
}

func NewServer(logger *slog.Logger, appEnv appenv.AppEnv, deps *Dependencies) (*httpServer, error) {
	maxUploadSize := deps.MaxUploadSize
	if maxUploadSize <= 0 {
		maxUploadSize = defaultMaxUploadSize
	}

	handlers, err := handlers.New(
		logger,
		&handlers.Dependencies{
			UploadedFileCacheRepo: deps.UploadedFileCacheRepo,
			DownloadsStorage:      deps.DownloadsStorage,
			Assets:                assets.NewAssets(deps.AssetsDir, deps.AssetFileCacheRepo),
			Usecases:              deps.Usecases,
			Template:              deps.Template,
			AppEnv:                appEnv,
			AppMode:               deps.AppMode,
			BaseURL:               deps.BaseURL,
			ShortLinkPrefix:       deps.ShortLinkPrefix,
			MaxUploadSize:         maxUploadSize,
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
		maxUploadSize:   maxUploadSize,
	}, nil
}

func (s *httpServer) ListenAndServe(ctx context.Context, port string, ready chan<- struct{}) error {
	addr := fmt.Sprintf(":%s", port)

	routes := routes.NewRoutes(s.middlewares, s.handlers, s.shortLinkPrefix)
	router := routes.Router()
	handler := nfasthttp.NewHandler(ctx, s.logger, s.appEnv, router.Handler)

	if s.metricsEnabled {
		handler = metrics.MiddlewareHandler(handler)
	}

	server := &fasthttp.Server{
		Handler: handler,

		// Required for large uploads: the body is read as a stream and never buffered in memory.
		StreamRequestBody: true,

		// Let handlers parse multipart bodies on demand.
		DisablePreParseMultipartForm: true,

		// Per-request overrides: only the upload route gets the long read timeout and the large
		// body limit; every other route keeps the strict server defaults below.
		HeaderReceived: func(header *fasthttp.RequestHeader) fasthttp.RequestConfig {
			if string(header.Method()) == fasthttp.MethodPost && isUploadPath(header.RequestURI()) {
				return fasthttp.RequestConfig{
					ReadTimeout:        time.Hour,
					MaxRequestBodySize: int(s.maxUploadSize) + uploadMultipartOverhead,
				}
			}

			return fasthttp.RequestConfig{} // zero values mean "use the server settings"
		},

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
		Concurrency: 4096,
		// Prevents a single IP from exhausting workers.
		// 0 disables the per-IP connection limit.
		MaxConnsPerIP: 0,
		// Periodically rotates long-lived connections
		// 0 disables the limit.
		MaxRequestsPerConn: 0,

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

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		if ctx.Err() == nil {
			s.logger.ErrorContext(ctx, "Failed to start server", "error", err)
		}
		return err
	}

	close(ready)

	err = server.Serve(listener)
	if err != nil && ctx.Err() == nil {
		s.logger.ErrorContext(ctx, fmt.Sprintf("error run server: %v", err))
		return err
	}

	log.Println("Server stopped")

	return nil
}

// isUploadPath reports whether the request URI points to the upload endpoint (query string ignored).
func isUploadPath(uri []byte) bool {
	if i := bytes.IndexByte(uri, '?'); i >= 0 {
		uri = uri[:i]
	}

	return string(uri) == httppaths.DownloaderUploadPath
}
