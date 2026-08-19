package memory

import "context"

type CacheRepository interface {
	Name() string
	Transaction(fn func(context.Context) error) error
}
