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

	// Cookies contains the HTTP Cookie header value.
	Cookies []string
}

type FetchOption func(*FetchOptions)

func DefaultFetchOptions() FetchOptions {
	return FetchOptions{
		LimitBytes: defaultLimitBytes,
		Timeout:    defaultTimeout,
	}
}

func ApplyFetchOptions(options *FetchOptions, opts ...FetchOption) {
	for _, opt := range opts {
		opt(options)
	}
}

func NewFetchOptions(opts ...FetchOption) FetchOptions {
	options := DefaultFetchOptions()

	ApplyFetchOptions(&options, opts...)

	return options
}

func FetchOptionWithOptions(options FetchOptions) FetchOption {
	return func(fetchOptions *FetchOptions) {
		*fetchOptions = options
	}
}

func FetchOptionsWithCookies(cookies []string) FetchOption {
	return func(options *FetchOptions) {
		options.Cookies = cookies
	}
}
