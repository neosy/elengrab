package sldownload

import (
	"context"
	"database/sql"
	"time"

	"github.com/neosy/elengrab/internal/repository/sqlite/lock"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

const (
	maxRetriesDefault = 5
	retryDelayDefault = 200 * time.Millisecond
)

type dbInterface interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type ctxTxKey struct {
}

type ctxTxLockedKey struct {
}

type writeCtxTxLocked struct {
	locked bool
}

type retryOptions struct {
	maxRetries int
	delay      time.Duration
}

func ctxWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	ctxWithValue := context.WithValue(ctx, ctxTxKey{}, tx)
	ctxWithValue = context.WithValue(ctxWithValue, ctxTxLockedKey{}, &writeCtxTxLocked{locked: true})
	return ctxWithValue
}

func txFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(ctxTxKey{}).(*sql.Tx)
	return tx, ok
}

func dbOrTx(ctx context.Context, db *sql.DB) dbInterface {
	var (
		dbtx dbInterface = db
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

func execContext(ctx context.Context, db *sql.DB, lock lock.WriteLocker, sqlQuery string, args []any, options retryOptions) error {
	var (
		err  error
		dbtx = dbOrTx(ctx, db)
	)

	for i := range options.maxRetries {
		if !txLocked(ctx) {
			lock.Lock()
		}
		_, err = dbtx.ExecContext(ctx, sqlQuery, args...)
		if !txLocked(ctx) {
			lock.Unlock()
		}
		if err == nil {
			break
		}

		if sqlError, ok := err.(*sqlite.Error); ok && sqlError.Code() == sqlite3.SQLITE_BUSY {
			if i+1 == options.maxRetries {
				return err
			}

			timer := time.NewTimer(options.delay)

			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
				// Let's continue
			}

			continue
		}

		return err
	}

	return nil
}
