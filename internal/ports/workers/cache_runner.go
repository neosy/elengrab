package pworkers

import "context"

type CacheRunner interface {
	CleanExpired(ctx context.Context) error
}
