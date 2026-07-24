package persistence

import "context"

type Transactional interface {
	Tx(ctx context.Context, fn func(ctx context.Context) error) error
	TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error
}
