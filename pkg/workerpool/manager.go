package workerpool

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxWorkersDefault uint = 3
	// Interval to wait before retrying when queue is empty or semaphore is full
	dispatchRetryInterval = 1000 * time.Millisecond
)

type Manager interface {
	Start(ctx context.Context) error
	Stop()
}

type ManagerOptions struct {
	MaxWorkers uint
}

type manager struct {
	logger     *slog.Logger
	workers    []Worker
	maxWorkers uint

	quit    chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
	running atomic.Bool

	semaphore chan struct{}
	jobStream chan Job
	jobQueue  jobQueue
}

func NewManager(logger *slog.Logger, options *ManagerOptions) *manager {
	var maxWorkers = maxWorkersDefault
	if options != nil && options.MaxWorkers > 0 {
		maxWorkers = options.MaxWorkers
	}

	return &manager{
		logger:     logger,
		workers:    make([]Worker, 0),
		maxWorkers: maxWorkers,
		quit:       make(chan struct{}),
		semaphore:  make(chan struct{}, maxWorkers),
		jobStream:  make(chan Job, 1),
		jobQueue:   *newJobQueue(100),
	}
}

func (m *manager) Start(ctx context.Context) error {
	if !m.running.CompareAndSwap(false, true) {
		err := errors.New("manager already running")
		m.logger.Error("Download manager already running")
		return err
	}

	for range m.maxWorkers {
		m.addWorker(ctx)
	}

	m.wg.Go(func() {
		defer m.running.Store(false)

		// Create a timer that fires immediately
		timer := time.NewTimer(0)

		timerReset := func() {
			timer.Reset(dispatchRetryInterval)
		}

		for {
			select {
			// Stop loop if quit signal is received
			case <-m.quit:
				return
			// Handle timer expiration
			case <-timer.C:
				func() {
					m.mu.Lock()
					defer m.mu.Unlock()

					// If queue is empty, reset timer and return
					if m.jobQueue.Len() == 0 {
						timerReset()
						return
					}

					// Try to acquire a semaphore slot to dispatch a job
					select {
					case m.semaphore <- struct{}{}:
						// Pop a job from the queue and send it to jobStream
						job, _ := m.jobQueue.Pop()
						m.jobStream <- job
					// If semaphore is full, reset timer to retry later
					default:
						timerReset()
					}
				}()
			}
		}
	})

	m.logger.Debug("Downloader manager running...")

	return nil
}

func (m *manager) Stop() {
	if !m.running.CompareAndSwap(true, false) {
		return
	}

	close(m.quit)
	m.wg.Wait()
}

func (m *manager) AddJob(job Job) {
	m.jobQueue.Push(job)
}

func (m *manager) addWorker(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	worker := newWorker(
		m.logger,
		uint(len(m.workers)),
		&m.wg,
	)
	m.workers = append(m.workers, worker)

	worker.Start(
		ctx,
		m.jobStream,
		m.quit,
		func() {
			<-m.semaphore
		},
	)
}
