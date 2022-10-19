package budget

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

type Repository interface {
	GetBudgetTx(tx *sql.Tx, ctx context.Context, userId int64) (*model.Budget, error)
	SaveBudget(ctx context.Context, userId int64, limit int64) error
}

type budgetDB struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &budgetDB{
		db: db,
	}
}

func (c *budgetDB) GetBudgetTx(tx *sql.Tx, ctx context.Context, userId int64) (*model.Budget, error) {
	const query = `
		select
			id,
			user_id,
			value,	
		    created_at
		from budget 
		where user_id = $1;
	`

	var budget model.Budget
	err := tx.QueryRowContext(ctx, query, userId).Scan(&budget.Id, &budget.UserId, &budget.Value, &budget.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &budget, err
}

func (c *budgetDB) SaveBudget(ctx context.Context, userId int64, limit int64) error {
	const query = `
		insert into budget(
			created_at,
		    updated_at,
			user_id,
			value
		) values (now(), now(), $1, $2)
		ON CONFLICT (user_id) 
		DO UPDATE SET value = $2, updated_at = now();
	`

	_, err := c.db.ExecContext(ctx, query, userId, limit)

	return err
}
