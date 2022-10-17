package db

import (
	"context"
	"database/sql"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
)

type selectedCurrencyDB struct {
	db *sql.DB
}

func New(db *sql.DB) selected_currency.Repository {
	return &selectedCurrencyDB{
		db: db,
	}
}

func (s *selectedCurrencyDB) SaveSelectedCurrency(ctx context.Context, code string, userId int64) error {
	const query = `
		insert into selected_currency(
			created_at,
			code,
			user_id
		) values (
			now(), $1, $2
		);
	`

	_, err := s.db.ExecContext(ctx, query, code, userId)

	return err
}

func (s *selectedCurrencyDB) GetSelectedCurrency(ctx context.Context, userId int64) (*model.SelectedCurrency, error) {
	const query = `
		select 
		       id,
		       code,
		       user_id,
		       created_at
		from selected_currency
		where user_id = $1
		order by created_at desc
		limit 1
	`

	var selectedCurrency model.SelectedCurrency

	err := s.db.QueryRowContext(ctx, query, userId).Scan(&selectedCurrency.Id, &selectedCurrency.Code, &selectedCurrency.UserId, &selectedCurrency.CreatedAt)

	return &selectedCurrency, err
}
