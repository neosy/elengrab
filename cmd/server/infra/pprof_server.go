package infra

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	iconfig "github.com/neosy/elengrab/internal/config"
)

func newPprofMux() *http.ServeMux {
	mux := http.NewServeMux()

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

	return mux
}

func StartPprofHTTPServer(ctx context.Context, cfg iconfig.PprofServerConfig) {
	if !cfg.Enable {
		return
	}

	srv := &http.Server{
		Addr:    net.JoinHostPort(cfg.Address, cfg.Port),
		Handler: newPprofMux(),
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Println("pprof shutdown error:", err)
		}
	}()

	go func() {
		log.Println("pprof listening on", srv.Addr)

		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Println(err)
		}
	}()
}
