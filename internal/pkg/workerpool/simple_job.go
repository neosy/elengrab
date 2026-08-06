package workerpool

import (
	"context"
	"errors"
)

// SimpleJob is a basic Job implementation that executes a provided handler function.
// It can be used when creating a dedicated Job type is unnecessary.
type SimpleJob struct {
	id      string
	name    string
	handler func(ctx context.Context, workerID uint64) error
}

// NewSimpleJob creates a new SimpleJob with the given identifier, name, and execution handler.
func NewSimpleJob(
	jobID string,
	jobName string,
	handler func(ctx context.Context, workerID uint64) error,
) Job {
	return &SimpleJob{
		id:      jobID,
		name:    jobName,
		handler: handler,
	}
}

// ID returns the unique identifier of the job.
func (j *SimpleJob) ID() string {
	return j.id
}

// Name returns the human-readable name of the job.
func (j *SimpleJob) Name() string {
	return j.name
}

// Execute runs the job handler with the provided context and worker ID.
func (j *SimpleJob) Execute(ctx context.Context, workerID uint64) error {
	if j.handler == nil {
		return errors.New("job handler is nil")
	}

	return j.handler(ctx, workerID)
}
