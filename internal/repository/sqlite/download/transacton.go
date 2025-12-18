package sldownload

import (
	"context"
	"database/sql"
	"sync"
	"time"

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

type ctxTxKey struct{}

type retryOptions struct {
	maxRetries int
	delay      time.Duration
}

func ctxWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ctxTxKey{}, tx)
}

func txFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(ctxTxKey{}).(*sql.Tx)
	return tx, ok
}

func dbOrTx(ctx context.Context, db *sql.DB) dbInterface {
	var dbtx dbInterface = db

	if tx, ok := txFromCtx(ctx); ok && tx != nil {
		dbtx = tx
	}

	return dbtx
}

func execContext(ctx context.Context, db *sql.DB, mu *sync.RWMutex, sqlQuery string, args []any, options retryOptions) error {
	var (
		err  error
		dbtx = dbOrTx(ctx, db)
	)

	for i := range options.maxRetries {
		mu.Lock()
		_, err = dbtx.ExecContext(ctx, sqlQuery, args...)
		mu.Unlock()
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
