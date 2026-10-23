package workerpool

import (
	"container/list"
	"context"
)

type Job interface {
	ID() string
	Name() string
	Execute(ctx context.Context, workerID uint64) error
}

type JobDispatcher interface {
	// IsTaskRunning reports whether the task with the specified job ID is currently running.
	// This method is thread-safe and can be called from any goroutine.
	IsRunningTask(jobID string) bool

	// HasJob reports whether a job exists in the queue or is currently running.
	// This method is thread-safe and can be called from any goroutine.
	HasJob(jobID string) bool

	// AddJob enqueues a new job for execution.
	// job: the job to be executed by a worker.
	// Thread-safe; can be called from any goroutine.
	AddJob(job Job) error

	// CancelJob removes a job with the given ID from the queue.
	// Returns true if the job was found and removed, false otherwise.
	// If the job is removed, it signals the dispatcher to wake up
	// and re-evaluate available jobs.
	CancelJob(jobID string) bool
}

type jobQueue struct {
	items *list.List
	index map[string]*list.Element
}

type task struct {
	ctx    context.Context
	cancel context.CancelFunc
	job    Job
}

// newJobQueue creates a new job queue with the specified capacity.
func newJobQueue(cap uint32) *jobQueue {
	return &jobQueue{
		items: list.New(),
		index: make(map[string]*list.Element, cap),
	}
}

// newTask creates a new task with the specified context and job.
func newTask(ctx context.Context, job Job) *task {
	ctx, cancel := context.WithCancel(ctx)

	return &task{
		ctx:    ctx,
		cancel: cancel,
		job:    job,
	}
}

// len returns the number of items in the queue.
//
// This method is not thread-safe. The caller must hold the mutex
// protecting the queue before calling Len.
func (q *jobQueue) len() int {
	return q.items.Len()
}

// contains reports whether a job with the given ID exists in the queue.
//
// This method is not thread-safe. The caller must hold the mutex
// protecting the queue before calling Contains.
func (q *jobQueue) contains(jobID string) bool {
	_, exists := q.index[jobID]
	return exists
}

// push adds a new item to the queue.
//
// This method is not thread-safe. The caller must hold the mutex
// protecting the queue before calling Push.
// If an item with the same ID already exists, it is not added and false is returned.
func (q *jobQueue) push(job Job) bool {
	id := job.ID()
	if _, exists := q.index[id]; exists {
		return false
	}

	q.index[id] = q.items.PushBack(job)

	return true
}

// pop removes and returns the first item in the queue.
//
// This method is not thread-safe. The caller must hold the mutex
// protecting the queue before calling Pop.
// If the queue is empty, it returns nil and false.
func (q *jobQueue) pop() (Job, bool) {
	el := q.items.Front()
	if el == nil {
		return nil, false
	}

	job := el.Value.(Job)

	if !q.remove(job.ID()) {
		return nil, false
	}

	return job, true
}

// remove removes an item from the queue by ID.
//
// This method is not thread-safe. The caller must hold the mutex
// protecting the queue before calling Remove.
// If no such item exists, it returns false.
func (q *jobQueue) remove(id string) bool {
	el, exists := q.index[id]
	if !exists {
		return false
	}

	q.items.Remove(el)
	delete(q.index, id)

	return true
}

// Cancel cancels the task context.
// It should be called when the task is no longer needed or has completed.
func (t *task) Cancel() {
	t.cancel()
}

// Close releases resources associated with the task.
func (t *task) Close() {
	t.Cancel()
}
