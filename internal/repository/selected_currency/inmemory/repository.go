package inmemory

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
)

var CurrencyNotFound = errors.New("currency not found")

type inmemory struct {
	selectedCurrency map[int64]selected_currency.SelectedCurrency
}

func New() selected_currency.Repository {
	return &inmemory{
		selectedCurrency: map[int64]selected_currency.SelectedCurrency{},
	}
}

func (i *inmemory) SaveSelectedCurrency(currency string, userId int64) error {
	i.selectedCurrency[userId] = selected_currency.SelectedCurrency{
		Currency: currency,
		UserId:   userId,
	}

	return nil
}

func (i *inmemory) GetSelectedCurrency(userId int64) (selected_currency.SelectedCurrency, error) {
	if cur, ok := i.selectedCurrency[userId]; ok {
		return cur, nil
	}

	return selected_currency.SelectedCurrency{}, CurrencyNotFound
}
