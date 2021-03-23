package utils

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"reflect"
)

type TxxFunction func(ctx context.Context, tx *sqlx.Tx) error

func In(where []string, args []interface{}, field string, arg interface{}) ([]string, []interface{}, error) {
	v := reflect.ValueOf(arg)
	if arg == nil {
		return where, args, nil
	}

	t := reflectx.Deref(v.Type())
	if t.Kind() != reflect.Slice {
		where = append(where, field+" = ?")
		args = append(args, arg)
		return where, args, nil
	}
	if v.Len() == 0 {
		return where, args, nil
	} else if v.Len() == 1 {
		slice := v.Slice(0, 1)
		where = append(where, field+" = ?")
		args = append(args, slice.Index(0).Interface())
		return where, args, nil
	}

	q, qArgs, err := sqlx.In(field+" in (?)", arg)
	if err != nil {
		return where, args, err
	}
	where = append(where, q)
	args = append(args, qArgs...)
	return where, args, nil
}

//InReadOnlyTransactionX runs callback in read transaction
func InReadOnlyTransactionX(db *sqlx.DB, ctx context.Context, f TxxFunction) error {
	tx, txErr := db.BeginTxx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	if txErr != nil {
		return newBeginErr(txErr)
	}

	defer func() {
		_ = tx.Commit()
	}()

	fErr := f(ctx, tx)
	if fErr != nil {
		return fErr
	}
	ctxErr := ctx.Err()
	if ctxErr != nil {
		return newContextErr(ctxErr)
	}

	return nil
}

//InWriteTransactionX runs callback in write transaction
func InWriteTransactionX(db *sqlx.DB, ctx context.Context, callback TxxFunction) (retErr error) {

	tx, txErr := db.BeginTxx(ctx, &sql.TxOptions{})
	if txErr != nil {
		retErr = newBeginErr(txErr)
		return
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				retErr = newRollbackErr(err)
			}
		}
	}()

	if fErr := callback(ctx, tx); fErr != nil {
		retErr = newCallbackErr(fErr)
		return
	}

	ctxErr := ctx.Err()
	if ctxErr != nil {
		retErr = newContextErr(ctxErr)
		return
	}

	if commitErr := tx.Commit(); commitErr != nil {
		retErr = newCommitErr(commitErr)
	} else {
		committed = true
	}

	return nil
}
