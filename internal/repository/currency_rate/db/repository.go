package db

import (
	"context"
	"database/sql"
	"errors"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
)

type currencyRateDB struct {
	db *sql.DB
}

func New(db *sql.DB) currency_rate.Repository {
	return &currencyRateDB{
		db: db,
	}
}

func (c *currencyRateDB) SaveRate(ctx context.Context, name string, value int64) error {
	const query = `
		insert into currency_rate(
			created_at,
			code, 
			value
		) values (
			now(), $1, $2
		);
	`

	_, err := c.db.ExecContext(ctx, query, name, value)

	return err
}

func (c *currencyRateDB) GetRateByCurrency(ctx context.Context, code string) (*model.CurrencyRate, error) {
	const query = `
		select 
		       id,
		       code,
		       value,
		       created_at
		from currency_rate
		where code = $1
		order by created_at desc
		limit 1
	`

	var rate model.CurrencyRate
	err := c.db.QueryRowContext(ctx, query, code).Scan(&rate.Id, &rate.Code, &rate.Value, &rate.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	return &rate, nil
}
