package selected_currency

import "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"

type Repository interface {
	SaveSelectedCurrency(currency string, userId int64) error
	GetSelectedCurrency(userId int64) (model.SelectedCurrency, error)
}
