package dbexec

import (
	"context"
	"database/sql"
)

// Tx executes fn within a transaction.
// If the context already contains a transaction, it is reused.
// Otherwise, a new transaction is created for the duration of fn.
func Tx(ctx context.Context, db *sql.DB, locker WriteLocker, fn func(ctx context.Context) error) error {
	var (
		tx       = txFromCtx(ctx)
		newCtx   = ctx
		isOpenTx = false
	)

	if tx == nil {
		locker.Lock()
		defer locker.Unlock()

		var err error
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		newCtx = ctxWithTx(ctx, tx)
		isOpenTx = true
	}

	if err := fn(newCtx); err != nil {
		tx.Rollback()
		return err
	}

	if isOpenTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// TxIndependent executes fn in a new independent transaction.
// If the context already contains a transaction, it is ignored and
// a new transaction is created for the duration of fn.
func TxIndependent(ctx context.Context, db *sql.DB, locker WriteLocker, fn func(ctx context.Context) error) error {
	locker.Lock()
	defer locker.Unlock()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	newCtx := ctxWithTx(ctx, tx)

	if err := fn(newCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
