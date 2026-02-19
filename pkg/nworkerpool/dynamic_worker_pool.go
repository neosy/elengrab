package nworkerpool

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type dynamicWorkerPool struct {
	logger  *slog.Logger
	workers []Worker

	poolSize int
	idleTime time.Duration

	quit chan struct{}
	wg   sync.WaitGroup

	mu   sync.Mutex
	cond *sync.Cond

	running atomic.Bool

	semaphore    chan struct{}
	taskStream   chan *task
	jobQueue     *jobQueue
	runningTasks map[string]*task
}

// NewDynamicWorkerPool creates a new dynamic worker pool.
//
// A dynamic worker pool spawns workers on demand: a new worker is
// created when a job arrives and there are no idle workers available.
// Workers automatically exit after being idle for a specified duration.
//
// logger: used for logging events and errors.
// options: optional configuration; nil uses defaults. Supports:
//   - PoolSize: maximum number of concurrent workers (defaults to defaultDynamicWorkerPoolSize)
//   - IdleTime: duration a worker stays alive when idle (defaults to defaultIdleTime)
//
// Returns an initialized *dynamicWorkerPool ready to Start().
func NewDynamicWorkerPool(logger *slog.Logger, opt *DynamicWorkerPoolOptions) *dynamicWorkerPool {
	var poolSize = defaultDynamicWorkerPoolSize
	if opt != nil && opt.PoolSize > 0 {
		poolSize = opt.PoolSize
	}

	var idleTime time.Duration = defaultIdleTime
	if opt != nil && opt.IdleTime > 0 {
		idleTime = opt.IdleTime
	}

	wp := &dynamicWorkerPool{
		logger:       logger,
		workers:      make([]Worker, 0, poolSize),
		poolSize:     poolSize,
		idleTime:     idleTime,
		quit:         make(chan struct{}),
		semaphore:    make(chan struct{}, poolSize),
		taskStream:   make(chan *task, 1),
		jobQueue:     newJobQueue(defaultJobQueueCap), // Initial but not final queue size
		runningTasks: make(map[string]*task, poolSize),
	}

	wp.cond = sync.NewCond(&wp.mu)

	return wp
}
