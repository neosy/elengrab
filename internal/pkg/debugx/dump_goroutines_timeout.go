package debugx

import (
	"bytes"
	"log"
	"log/slog"
	"runtime/pprof"
	"time"
)

// DumpGoroutinesIfTimeoutWithLogger writes a goroutine dump to the logger if
// the operation does not complete within the specified timeout. The returned
// function must be deferred by the caller to stop the timeout watcher.
func DumpGoroutinesIfTimeoutWithLogger(
	logger *slog.Logger,
	operation string,
	timeout time.Duration,
) func() {
	if logger == nil {
		logger = slog.Default()
	}

	done := make(chan struct{})

	go func() {
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		select {
		case <-done:
			return

		case <-timer.C:
			var buf bytes.Buffer

			if err := pprof.Lookup("goroutine").WriteTo(&buf, 2); err != nil {
				logger.Error(
					"Failed to dump goroutines",
					slog.Any("error", err),
				)
				return
			}

			logger.Warn(
				"Slow operation detected",
				slog.String("operation", operation),
				slog.Duration("timeout", timeout),
				slog.String("goroutines", buf.String()),
			)
		}
	}()

	return func() {
		close(done)
	}
}

// DumpGoroutinesIfTimeout writes a goroutine dump to the default logger if
// the operation does not complete within the specified timeout. The returned
// function must be deferred by the caller to stop the timeout watcher.
func DumpGoroutinesIfTimeout(
	operation string,
	timeout time.Duration,
) func() {
	done := make(chan struct{})

	go func() {
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		select {
		case <-done:
			return

		case <-timer.C:
			var buf bytes.Buffer

			if err := pprof.Lookup("goroutine").WriteTo(&buf, 2); err != nil {
				log.Printf("failed to dump goroutines: %v", err)
				return
			}

			log.Printf(
				"Slow operation detected: %s (%v)\n%s",
				operation,
				timeout,
				buf.String(),
			)
		}
	}()

	return func() {
		close(done)
	}
}
