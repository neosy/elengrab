package idto

type ExecutorOption func(*ExecutorOptions)

type ExecutorOptions struct {
	EnsureCache    bool
	CookieFilePath string
}

func NewExecutorOptionsDefault() *ExecutorOptions {
	return &ExecutorOptions{
		EnsureCache:    true,
		CookieFilePath: "",
	}
}

func NewExecutorOptions(opts ...ExecutorOption) *ExecutorOptions {
	options := NewExecutorOptionsDefault()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithEnsureCache(enabled bool) ExecutorOption {
	return func(o *ExecutorOptions) {
		o.EnsureCache = enabled
	}
}

func WithUseCookies(filePath string) ExecutorOption {
	return func(o *ExecutorOptions) {
		o.CookieFilePath = filePath
	}
}
