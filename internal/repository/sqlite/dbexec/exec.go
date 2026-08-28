package dbexec

import (
	"context"
	"database/sql"
	"time"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type DBExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func ExecContext(ctx context.Context, dbEntry dbEntry, sqlQuery string, args []any, options RetryOptions) error {
	var (
		err  error
		dbtx = Resolve(ctx, dbEntry)
	)

	for i := range options.MaxRetries {
		if !txLocked(ctx) {
			dbEntry.Locker().Lock()
		}
		_, err = dbtx.ExecContext(ctx, sqlQuery, args...)
		if !txLocked(ctx) {
			dbEntry.Locker().Unlock()
		}
		if err == nil {
			break
		}

		if sqlError, ok := err.(*sqlite.Error); ok && sqlError.Code() == sqlite3.SQLITE_BUSY {
			if i+1 == options.MaxRetries {
				return err
			}

			timer := time.NewTimer(options.Delay)

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

func ExecRetry(
	ctx context.Context,
	opts RetryOptions,
	fn func() error,
) error {
	for i := range opts.MaxRetries {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := fn()
		if err == nil {
			return nil
		}

		if sqlError, ok := err.(*sqlite.Error); ok && sqlError.Code() == sqlite3.SQLITE_BUSY {
			if i+1 == opts.MaxRetries {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(opts.Delay):
				// Let's continue
			}

			continue
		}

		return err
	}

	return nil
}
