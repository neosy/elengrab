package memsimple

import "context"

// transactionContextKey identifies a context carrying an active repository transaction.
type transactionContextKey struct{}

// withTransactionContext adds a transaction marker to the context.
func withTransactionContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, transactionContextKey{}, struct{}{})
}

// isTransactionContext checks whether the context contains a transaction marker.
func isTransactionContext(ctx context.Context) bool {
	return ctx != nil && ctx.Value(transactionContextKey{}) != nil
}
