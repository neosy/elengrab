package pworkers

import "context"

type AuthWebStartupRunner interface {
	Startup(ctx context.Context) error
}
