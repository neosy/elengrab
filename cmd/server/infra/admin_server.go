package infra

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	iconfig "github.com/neosy/elengrab/internal/config"
)

func newAdminServerMux(cfg iconfig.AdminServerDebugConfig) *http.ServeMux {
	mux := http.NewServeMux()

	// pprof
	if cfg.EnablePprof {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		// explicit profiles (optional but clearer)
		mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
		mux.Handle("/debug/pprof/block", pprof.Handler("block"))
		mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	}

	// metrics
	if cfg.EnableMetrics {
		mux.Handle("/metrics", promhttp.Handler())
	}

	// health
	if cfg.EnableHealth {
		mux.HandleFunc("/healthz", healthHandler)
	}

	return mux
}

func StartAdminHTTPServer(ctx context.Context, cfg iconfig.AdminServerConfig) {
	if !cfg.Enable {
		return
	}

	srv := &http.Server{
		Addr:    net.JoinHostPort(cfg.Address, cfg.Port),
		Handler: newAdminServerMux(cfg.DebugConfig),
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Println("Admin server shutdown error:", err)
		}
	}()

	go func() {
		log.Println("Admin server listening on", srv.Addr)

		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Println(err)
		}
	}()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
