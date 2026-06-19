// manager_test.go
package nworkerpool

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
)

// noopLogger returns a logger that discards all output except errors.
func noopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

// ---------------------------------------------------------------------
// testJob – simple job that records execution
// ---------------------------------------------------------------------
type testJob struct {
	id       int
	executed atomic.Bool
	done     chan struct{}
	workerID uint64
	sleep    time.Duration
}

func (j *testJob) ID() string {
	return strconv.Itoa(j.id)
}

func (j *testJob) Name() string {
	return "Test"
}

func (j *testJob) Execute(ctx context.Context, workerID uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	j.workerID = workerID
	j.executed.Store(true)

	if j.sleep > 0 {
		select {
		case <-ctx.Done():
			close(j.done)
			return ctx.Err()
		case <-time.After(j.sleep):
		}
	}

	close(j.done)
	return nil
}

// ---------------------------------------------------------------------
// concurrentJobWrapper – wraps any Job to track concurrency
// ---------------------------------------------------------------------
type concurrentJobWrapper struct {
	nworkerpool.Job
	running *atomic.Int32
	maxSeen *atomic.Int32
	wg      *sync.WaitGroup
}

func (w *concurrentJobWrapper) Execute(ctx context.Context, workerID uint64) error {
	count := w.running.Add(1)
	if count > w.maxSeen.Load() {
		w.maxSeen.Store(count)
	}
	defer w.running.Add(-1)
	defer w.wg.Done()
	return w.Job.Execute(ctx, workerID)
}

// ---------------------------------------------------------------------
// TestManager_StartStop
// ---------------------------------------------------------------------
func TestManager_StartStop(t *testing.T) {
	m := nworkerpool.NewWorkerPool("", nworkerpool.WithLogger(noopLogger()))
	if err := m.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	m.Stop()
	m.Stop() // idempotent
}

// ---------------------------------------------------------------------
// TestManager_DoubleStart
// ---------------------------------------------------------------------
func TestManager_DoubleStart(t *testing.T) {
	m := nworkerpool.NewWorkerPool("", nworkerpool.WithLogger(noopLogger()))
	_ = m.Start(context.Background())
	if err := m.Start(context.Background()); err == nil {
		t.Fatal("expected error on double start")
	}
}

// ---------------------------------------------------------------------
// TestManager_JobExecution
// ---------------------------------------------------------------------
func TestManager_JobExecution(t *testing.T) {
	const workers = 2
	m := nworkerpool.NewWorkerPool(
		"",
		nworkerpool.WithLogger(noopLogger()),
		nworkerpool.WithMaxWorkers(workers),
	)
	if err := m.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer m.Stop()

	const jobs = 5
	var wg sync.WaitGroup
	wg.Add(jobs)

	for i := range jobs {
		j := &testJob{id: i, done: make(chan struct{})}
		if !m.AddJob(j) {
			continue
		}

		go func(job *testJob) {
			defer wg.Done()
			<-job.done
			if !job.executed.Load() {
				t.Errorf("job %d not executed", job.id)
			}
		}(j)
	}

	if !waitWithTimeout(&wg, 2*time.Second) {
		t.Fatal("timeout")
	}
}

// ---------------------------------------------------------------------
// TestManager_MaxWorkersLimit
// ---------------------------------------------------------------------
func TestManager_MaxWorkersLimit(t *testing.T) {
	const maxWorkers = 3
	m := nworkerpool.NewWorkerPool(
		"",
		nworkerpool.WithLogger(noopLogger()),
		nworkerpool.WithMaxWorkers(maxWorkers),
	)
	if err := m.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer m.Stop()

	var (
		running atomic.Int32
		maxSeen atomic.Int32
		wg      sync.WaitGroup
		total   = maxWorkers*2 + 1
	)

	for i := range total {
		base := &testJob{
			id:    i,
			done:  make(chan struct{}),
			sleep: 200 * time.Millisecond,
		}
		wrapped := &concurrentJobWrapper{
			Job:     base,
			running: &running,
			maxSeen: &maxSeen,
			wg:      &wg,
		}
		if m.AddJob(wrapped) {
			wg.Add(1)
		}
	}

	if !waitWithTimeout(&wg, 5*time.Second) {
		t.Fatal("timeout")
	}

	if maxSeen.Load() > maxWorkers {
		t.Errorf("concurrency exceeded: %d > %d", maxSeen.Load(), maxWorkers)
	}
}

// ---------------------------------------------------------------------
// TestManager_ContextCancellation
// ---------------------------------------------------------------------
func TestManager_ContextCancellation(t *testing.T) {
	m := nworkerpool.NewWorkerPool(
		"",
		nworkerpool.WithLogger(noopLogger()),
		nworkerpool.WithMaxWorkers(1),
	)

	// Создаём контекст
	ctx, cancel := context.WithCancel(context.Background())
	if err := m.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer m.Stop()

	j := &testJob{
		done:  make(chan struct{}),
		sleep: 10 * time.Second,
	}

	m.AddJob(j)

	// Даём воркеру начать
	time.Sleep(50 * time.Millisecond)

	// ОТМЕНЯЕМ
	cancel()

	// Ждём завершения job'а
	select {
	case <-j.done:
		// Успешно: job завершился из-за отмены
	case <-time.After(1 * time.Second):
		t.Fatal("job did not stop after cancel")
	}
}

// ---------------------------------------------------------------------
// TestManager_ImmediateStop
// ---------------------------------------------------------------------
func TestManager_ImmediateStop(t *testing.T) {
	m := nworkerpool.NewWorkerPool("", nworkerpool.WithLogger(noopLogger()))
	if err := m.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
	start := time.Now()
	m.Stop()
	if time.Since(start) > 100*time.Millisecond {
		t.Fatal("Stop too slow")
	}
}

// ---------------------------------------------------------------------
// TestManager_ZeroWorkers
// ---------------------------------------------------------------------
func TestManager_ZeroWorkers(t *testing.T) {
	m := nworkerpool.NewWorkerPool(
		"",
		nworkerpool.WithLogger(noopLogger()),
		nworkerpool.WithMaxWorkers(0),
	)
	if err := m.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer m.Stop()

	j := &testJob{done: make(chan struct{})}
	m.AddJob(j)

	select {
	case <-j.done:
		if m.PoolSize() != 3 {
			t.Fatal("Default worker count is not 3")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Job was not run")
	}
}

// ---------------------------------------------------------------------
// TestManager_Race
// ---------------------------------------------------------------------
func TestManager_Race(t *testing.T) {
	m := nworkerpool.NewWorkerPool(
		"",
		nworkerpool.WithLogger(noopLogger()),
		nworkerpool.WithMaxWorkers(5),
	)
	if err := m.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer m.Stop()

	const producers, jobsPer = 10, 20
	var wg sync.WaitGroup

	for range producers {
		wg.Go(func() {
			for range jobsPer {
				j := &testJob{done: make(chan struct{})}
				if m.AddJob(j) {
					<-j.done
				}
			}
		})
	}
	wg.Wait()
}

// ---------------------------------------------------------------------
// waitWithTimeout waits for WaitGroup with timeout
// ---------------------------------------------------------------------
func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}

// ---------------------------------------------------------------------
// TestManager_CancelJob verifies that CancelJob correctly removes jobs
// from the queue, returns appropriate results, and does not affect
// other jobs that are still in the queue.
// ---------------------------------------------------------------------
func TestManager_CancelJob(t *testing.T) {
	m := nworkerpool.NewWorkerPool(
		"",
		nworkerpool.WithLogger(noopLogger()),
		nworkerpool.WithMaxWorkers(1),
	)
	if err := m.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer m.Stop()

	// Add a job with sleep to ensure it does not start immediately
	job1 := &testJob{id: 1, done: make(chan struct{}), sleep: 500 * time.Millisecond}
	if !m.AddJob(job1) {
		t.Fatal("failed to add job1")
	}

	// Cancel the job before it executes; should return true
	if !m.CancelJob(job1.ID()) {
		t.Error("expected CancelJob to return true for existing job")
	}

	// Cancel the same job again; should return false
	if m.CancelJob(job1.ID()) {
		t.Error("expected CancelJob to return false for already removed job")
	}

	// Cancel a nonexistent job; should return false
	if m.CancelJob("nonexistent") {
		t.Error("expected CancelJob to return false for nonexistent job")
	}

	// Add another job and ensure it executes normally
	job2 := &testJob{id: 2, done: make(chan struct{})}
	if !m.AddJob(job2) {
		t.Fatal("failed to add job2")
	}

	select {
	case <-job2.done:
		// success: job executed
	case <-time.After(1 * time.Second):
		t.Fatal("job2 did not complete")
	}
}
