package database

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	"go.uber.org/zap"
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
			logger.Error("error in recover")
		}
	}()

	ok, errCb := callback(tx, ctx)
	if errCb != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error("error rolling back a transaction", zap.Error(rollbackErr))
		}

		return errors.WithStack(errCb)
	}

	if ok {
		return errors.WithStack(tx.Commit())
	}

	return errors.WithStack(tx.Rollback())
}
