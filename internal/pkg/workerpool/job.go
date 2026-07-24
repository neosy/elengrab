package nworkerpool

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
	AddJob(job Job) bool
	CancelJob(jobID string) bool
}

type jobQueue struct {
	items *list.List
	index map[string]*list.Element
}

type task struct {
	ctx    context.Context
	Cancel context.CancelFunc
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
		Cancel: cancel,
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
