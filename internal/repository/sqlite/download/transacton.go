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

func execContext(ctx context.Context, db *sql.DB, sqlQuery string, args []interface{}, options retryOptions) error {
	var err error
	if tx, ok := txFromCtx(ctx); ok && tx != nil {
		for i := range options.maxRetries {
			_, err = tx.ExecContext(ctx, sqlQuery, args...)
			if err != nil {
				if sqlError, ok := err.(*sqlite.Error); ok && sqlError.Code() == sqlite3.SQLITE_BUSY {
					if i == options.maxRetries-1 {
						return err
					}
					time.Sleep(options.delay)
					continue
				} else {
					return err
				}
			}
			break
		}
	} else {
		for i := range options.maxRetries {
			_, err = db.ExecContext(ctx, sqlQuery, args...)
			if err != nil {
				if sqlError, ok := err.(*sqlite.Error); ok && sqlError.Code() == sqlite3.SQLITE_BUSY {
					if i == options.maxRetries-1 {
						return err
					}
					time.Sleep(options.delay)
					continue
				} else {
					return err
				}
			}
			break
		}
	}

	return nil
}
