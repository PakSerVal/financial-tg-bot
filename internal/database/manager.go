package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/pkg/errors"
)

type SqlManager interface {
	InTransaction(ctx context.Context, callback func(tx *sql.Tx, ctx context.Context) (bool, error)) error
}

type sqlManager struct {
	db *sql.DB
}

func NewSqlManager(db *sql.DB) SqlManager {
	return &sqlManager{db: db}
}

func (m *sqlManager) InTransaction(ctx context.Context, callback func(tx *sql.Tx, ctx context.Context) (bool, error)) error {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if p := recover(); p != nil {
			err = tx.Rollback()
			panic(p)
		}
	}()

	ok, errCb := callback(tx, ctx)
	if errCb != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Println("error rolling back a transaction: ", rollbackErr)
		}

		return errors.WithStack(errCb)
	}

	if ok {
		return errors.WithStack(tx.Commit())
	}

	return errors.WithStack(tx.Rollback())
}
