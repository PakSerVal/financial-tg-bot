package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

type spendDB struct {
	db *sql.DB
}

func New(db *sql.DB) spend.Repository {
	return &spendDB{
		db: db,
	}
}

func (s *spendDB) SaveTx(tx *sql.Tx, ctx context.Context, price int64, category string, userId int64) error {
	const query = `
		insert into spend(
			created_at,
			price, 
		    category,
			user_id
		) values (
			now(), $1, $2, $3
		);
	`

	_, err := tx.ExecContext(ctx, query, price, category, userId)

	return err
}

func (s *spendDB) GetByTimeSince(ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error) {
	const query = `
		select 
		       id,
		       price,
		       category,
		       user_id,
		       created_at
		from spend
		where user_id = $1 and created_at > $2
	`

	rows, err := s.db.QueryContext(ctx, query, userId, timeSince)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	var spends []model.Spend
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		var spendModel model.Spend
		err = rows.Scan(&spendModel.Id, &spendModel.Price, &spendModel.Category, &spendModel.UserId, &spendModel.CreatedAt)
		if err != nil {
			return nil, err
		}

		spends = append(spends, spendModel)
	}

	return spends, err
}

func (s *spendDB) GetByTimeSinceTx(tx *sql.Tx, ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error) {
	const query = `
		select 
		       id,
		       price,
		       category,
		       user_id,
		       created_at
		from spend
		where user_id = $1 and created_at > $2
	`

	rows, err := tx.QueryContext(ctx, query, userId, timeSince)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	var spends []model.Spend
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		var spendModel model.Spend
		err = rows.Scan(&spendModel.Id, &spendModel.Price, &spendModel.Category, &spendModel.UserId, &spendModel.CreatedAt)
		if err != nil {
			return nil, err
		}

		spends = append(spends, spendModel)
	}

	return spends, err
}
