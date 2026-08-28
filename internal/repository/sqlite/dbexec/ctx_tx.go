package dbexec

import (
	"context"
	"database/sql"
)

type ctxTxKey struct {
	dbName string
}

type ctxTxLockedKey struct {
}

type writeCtxTxLocked struct {
	locked bool
}

func ctxWithTx(ctx context.Context, dbName string, tx *sql.Tx) context.Context {
	ctxWithValue := context.WithValue(ctx, ctxTxKey{dbName}, tx)
	ctxWithValue = context.WithValue(ctxWithValue, ctxTxLockedKey{}, &writeCtxTxLocked{locked: true})
	return ctxWithValue
}

func txFromCtx(ctx context.Context, dbName string) *sql.Tx {
	if tx, ok := ctx.Value(ctxTxKey{dbName}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func txLocked(ctx context.Context) bool {
	txLocked, ok := ctx.Value(ctxTxLockedKey{}).(*writeCtxTxLocked)
	return ok && txLocked.locked
}

func Resolve(ctx context.Context, dbEntry dbEntry) DBExecutor {
	if tx := txFromCtx(ctx, dbEntry.DBName()); tx != nil {
		return tx
	}

	return dbEntry.DB()
}
