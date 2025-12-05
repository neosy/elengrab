package sldownload

import (
	"context"
	"database/sql"
	"time"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

const (
	maxRetriesDefault = 5
	retryDelayDefault = 200 * time.Millisecond
)

type ctxTxKey struct{}

func ctxWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ctxTxKey{}, tx)
}

func txFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(ctxTxKey{}).(*sql.Tx)
	return tx, ok
}

type retryOptions struct {
	maxRetries int
	delay      time.Duration
}

func execContext(ctx context.Context, db *sql.DB, sqlQuery string, args []any, options retryOptions) error {
	var (
		err  error
		dbtx interface {
			ExecContext(context.Context, string, ...any) (sql.Result, error)
		} = db
	)

	if tx, ok := txFromCtx(ctx); ok && tx != nil {
		dbtx = tx
	}

	for i := range options.maxRetries {
		_, err = dbtx.ExecContext(ctx, sqlQuery, args...)
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
		} else {
			return err
		}
	}

	return nil
}
