package downloader

import dauth "github.com/neosy/elengrab/internal/domain/auth"

type callOptions struct {
	authCtx *dauth.UserContext
}

type callOption func(*callOptions)

func defaultCallOptions() callOptions {
	return callOptions{}
}

func buildCallOptions(opts ...callOption) callOptions {
	callOptions := defaultCallOptions()

	for _, opt := range opts {
		opt(&callOptions)
	}

	return callOptions
}

func withAuth(authCtx dauth.UserContext) callOption {
	return func(co *callOptions) {
		co.authCtx = &authCtx
	}
}
