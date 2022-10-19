package inmemory

import (
	"context"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/err_msg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
)

type inmemory struct {
	selectedCurrency map[int64]model.SelectedCurrency
}

func New() selected_currency.Repository {
	return &inmemory{
		selectedCurrency: map[int64]model.SelectedCurrency{},
	}
}

func (i *inmemory) SaveSelectedCurrency(ctx context.Context, currency string, userId int64) error {
	i.selectedCurrency[userId] = model.SelectedCurrency{
		Code:   currency,
		UserId: userId,
	}

	return nil
}

func (i *inmemory) GetSelectedCurrency(ctx context.Context, userId int64) (*model.SelectedCurrency, error) {
	if cur, ok := i.selectedCurrency[userId]; ok {
		return &cur, nil
	}

	return nil, err_msg.CurrencyNotFound
}
