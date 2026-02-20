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

// Len returns the number of items in the queue.
func (q *jobQueue) Len() int {
	return q.items.Len()
}

// Push adds a new item to the queue.
// If an item with the same ID already exists, it is not added and false is returned.
func (q *jobQueue) Push(job Job) bool {
	id := job.ID()
	if _, exists := q.index[id]; exists {
		return false
	}

	q.index[id] = q.items.PushBack(job)

	return true
}

// Pop removes and returns the first item in the queue.
// If the queue is empty, it returns nil and false.
func (q *jobQueue) Pop() (Job, bool) {
	el := q.items.Front()
	if el == nil {
		return nil, false
	}

	job := el.Value.(Job)

	if !q.Remove(job.ID()) {
		return nil, false
	}

	return job, true
}

// Remove removes an item from the queue by ID.
// If no such item exists, it returns false.
func (q *jobQueue) Remove(id string) bool {
	el, exists := q.index[id]
	if !exists {
		return false
	}

	q.items.Remove(el)
	delete(q.index, id)

	return true
}
