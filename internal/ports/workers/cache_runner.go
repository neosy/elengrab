package pworkers

import "context"

type CacheRunner interface {
	Name() string
	CleanExpired(ctx context.Context) error
}
