package nworkers

// JobExecuteOptions holds configuration options for executing a job.
type JobExecuteOptions struct {
	measureElapsed bool
}

// JobExecuteOption is a function type that modifies JobExecuteOptions.
type JobExecuteOption func(*JobExecuteOptions)

// DefaultJobExecuteOptions returns the default JobExecuteOptions.
func DefaultJobExecuteOptions() JobExecuteOptions {
	return JobExecuteOptions{
		measureElapsed: true,
	}
}

// ApplyJobExecuteOptions applies the provided JobExecuteOption functions to the given JobExecuteOptions.
func ApplyJobExecuteOptions(options *JobExecuteOptions, opts ...JobExecuteOption) {
	for _, opt := range opts {
		opt(options)
	}
}

// NewJobExecuteOptions creates a new JobExecuteOptions instance with the provided options applied.
func NewJobExecuteOptions(opts ...JobExecuteOption) JobExecuteOptions {
	options := DefaultJobExecuteOptions()

	ApplyJobExecuteOptions(&options, opts...)

	return options
}

// WithMeasureElapsed returns a JobExecuteOption that sets the measureElapsed field of JobExecuteOptions.
func WithMeasureElapsed(measureElapsed bool) JobExecuteOption {
	return func(opts *JobExecuteOptions) {
		opts.measureElapsed = measureElapsed
	}
}
