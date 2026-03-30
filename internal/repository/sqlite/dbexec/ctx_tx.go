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

func ctxWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	ctxWithValue := context.WithValue(ctx, ctxTxKey{}, tx)
	ctxWithValue = context.WithValue(ctxWithValue, ctxTxLockedKey{}, &writeCtxTxLocked{locked: true})
	return ctxWithValue
}

func txFromCtx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(ctxTxKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func txLocked(ctx context.Context) bool {
	txLocked, ok := ctx.Value(ctxTxLockedKey{}).(*writeCtxTxLocked)
	return ok && txLocked.locked
}

func Resolve(ctx context.Context, db *sql.DB) DBExecutor {
	if tx := txFromCtx(ctx); tx != nil {
		return tx
	}

	return db
}
