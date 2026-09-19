package types

import "time"

const (
	defaultLimitBytes = 4096 << 10 // 4 MiB
	defaultTimeout    = 10 * time.Second
)

type FetchOptions struct {
	// LimitBytes limits the maximum image size. Zero uses the default limit.
	LimitBytes int64

	// Timeout limits the image request duration. Zero uses the default timeout.
	Timeout time.Duration
}

func DefaultFetchOptions() FetchOptions {
	return FetchOptions{
		LimitBytes: defaultLimitBytes,
		Timeout:    defaultTimeout,
	}
}

func NewFetchOptions(options FetchOptions) FetchOptions {
	fetchOptions := DefaultFetchOptions()

	if options.LimitBytes != 0 {
		fetchOptions.LimitBytes = options.LimitBytes
	}

	if options.Timeout != 0 {
		fetchOptions.Timeout = options.Timeout
	}

	return fetchOptions
}
