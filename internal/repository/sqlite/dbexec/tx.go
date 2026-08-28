package dbexec

import (
	"context"
)

// Tx executes fn within a transaction.
// If the context already contains a transaction, it is reused.
// Otherwise, a new transaction is created for the duration of fn.
func Tx(ctx context.Context, dbEntry dbEntry, fn func(ctx context.Context) error) error {
	var (
		tx       = txFromCtx(ctx, dbEntry.DBName())
		newCtx   = ctx
		isOpenTx = false
	)

	if tx == nil {
		dbEntry.Locker().Lock()
		defer dbEntry.Locker().Unlock()

		var err error
		tx, err = dbEntry.DB().BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		newCtx = ctxWithTx(ctx, dbEntry.DBName(), tx)
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
func TxIndependent(ctx context.Context, dbEntry dbEntry, fn func(ctx context.Context) error) error {
	dbEntry.Locker().Lock()
	defer dbEntry.Locker().Unlock()

	tx, err := dbEntry.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	newCtx := ctxWithTx(ctx, dbEntry.DBName(), tx)

	if err := fn(newCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
