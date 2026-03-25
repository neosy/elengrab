package dbexec

import (
	"context"
	"database/sql"
)

func Tx(ctx context.Context, db *sql.DB, locker WriteLocker, fn func(ctx context.Context) error) error {
	tx, _ := TxFromCtx(ctx)

	var isOpenTx bool
	var newCtx = ctx
	if tx == nil {
		locker.Lock()
		defer locker.Unlock()

		var err error
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		newCtx = CtxWithTx(ctx, tx)
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
