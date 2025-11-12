package workerpool

import (
	"context"
	"sync"
)

type Job interface {
	Execute(ctx context.Context, workerId uint) error
}

type JobDispatcher interface {
	AddJob(job Job)
}

type jobQueue struct {
	items []Job
	mu    sync.Mutex
}

func newJobQueue(cap int) *jobQueue {
	return &jobQueue{
		items: make([]Job, 0, cap),
	}
}

func (q *jobQueue) Len() int {
	return len(q.items)
}

func (q *jobQueue) Push(job Job) {
	q.mu.Lock()
	q.items = append(q.items, job)
	q.mu.Unlock()
}

func (q *jobQueue) Pop() (Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return nil, false
	}

	job := q.items[0]
	q.items = q.items[1:]

	return job, true
}
