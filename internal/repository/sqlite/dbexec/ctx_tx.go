package dbexec

import (
	"context"
	"database/sql"
)

type ctxTxKey struct {
}

type ctxTxLockedKey struct {
}

type writeCtxTxLocked struct {
	locked bool
}

func CtxWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	ctxWithValue := context.WithValue(ctx, ctxTxKey{}, tx)
	ctxWithValue = context.WithValue(ctxWithValue, ctxTxLockedKey{}, &writeCtxTxLocked{locked: true})
	return ctxWithValue
}

func txFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(ctxTxKey{}).(*sql.Tx)
	return tx, ok
}

func Resolve(ctx context.Context, db *sql.DB) DBExecutor {
	var (
		dbtx DBExecutor = db
	)

	if tx, ok := txFromCtx(ctx); ok && tx != nil {
		dbtx = tx
	}

	return dbtx
}

func txLocked(ctx context.Context) bool {
	txLocked, ok := ctx.Value(ctxTxLockedKey{}).(*writeCtxTxLocked)
	return ok && txLocked.locked
}
