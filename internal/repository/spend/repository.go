package spend

import (
	"context"
	"database/sql"
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

type Repository interface {
	GetByTimeSince(ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error)
	SaveTx(tx *sql.Tx, ctx context.Context, sum int64, category string, userId int64) error
	GetByTimeSinceTx(tx *sql.Tx, ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error)
}
