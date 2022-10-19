package selected_currency

import (
	"context"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

type Repository interface {
	SaveSelectedCurrency(ctx context.Context, currency string, userId int64) error
	GetSelectedCurrency(ctx context.Context, userId int64) (*model.SelectedCurrency, error)
}
