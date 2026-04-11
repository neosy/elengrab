package nworkerpool

import (
	"context"
	"log/slog"
	"testing"
	"time"

	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
)

func TestManagerStopWhenIdle(t *testing.T) {
	logger := slog.Default()
	m := nworkerpool.NewWorkerPool(
		nworkerpool.WorkerPoolOptionLogger(logger),
		nworkerpool.WorkerPoolOptionMaxWorkers(1),
	)
	ctx := context.Background()

	if err := m.Start(ctx); err != nil {
		t.Fatal(err)
	}

	// Даём диспетчеру уснуть
	time.Sleep(100 * time.Millisecond)

	start := time.Now()
	m.Stop()
	duration := time.Since(start)

	if duration > 50*time.Millisecond {
		t.Errorf("Stop() too slow: %v", duration)
	}
}
