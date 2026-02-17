package executor

type Option func(*options)

type options struct {
	ensureCache bool
	useCookies  bool
}

func newDefaultOptions() *options {
	return &options{
		ensureCache: true,
		useCookies:  false,
	}
}

func WithEnsureCache(enabled bool) Option {
	return func(o *options) {
		o.ensureCache = enabled
	}
}

func WithUseCookies(enabled bool) Option {
	return func(o *options) {
		o.useCookies = enabled
	}
}
