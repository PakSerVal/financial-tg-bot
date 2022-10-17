package currency_rate

import (
	"context"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

type Repository interface {
	SaveRate(ctx context.Context, name string, value int64) error
	GetRateByCurrency(ctx context.Context, currency string) (*model.CurrencyRate, error)
}
