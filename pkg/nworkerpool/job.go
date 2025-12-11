package nworkerpool

import (
	"container/list"
	"context"
)

type Job interface {
	ID() string
	Name() string
	Execute(ctx context.Context, workerID uint) error
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

func newJobQueue(cap int) *jobQueue {
	return &jobQueue{
		items: list.New(),
		index: make(map[string]*list.Element, cap),
	}
}

func newTask(ctx context.Context, job Job) *task {
	ctx, cancel := context.WithCancel(ctx)

	return &task{
		ctx:    ctx,
		Cancel: cancel,
		job:    job,
	}
}

func (q *jobQueue) Len() int {
	return q.items.Len()
}

func (q *jobQueue) Push(job Job) bool {
	id := job.ID()
	if _, exists := q.index[id]; exists {
		return false
	}

	q.index[id] = q.items.PushBack(job)

	return true
}

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

func (q *jobQueue) Remove(id string) bool {
	el, exists := q.index[id]
	if !exists {
		return false
	}

	q.items.Remove(el)
	delete(q.index, id)

	return true
}
