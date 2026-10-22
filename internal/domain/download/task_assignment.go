package ddownload

import "uuid"

// TaskAssignment describes which worker and job are currently processing a task.
// Both fields are nil while the task is not being processed.
// It is meant to be embedded into task structs.
type TaskAssignment struct {
	// ID of the worker currently processing the task
	WorkerID *uint64

	// ID of the job currently processing the task
	JobID *uuid.UUID
}
