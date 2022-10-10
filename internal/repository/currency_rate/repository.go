package currency_rate

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

var ErrCurrencyRateNotFound = errors.New("repo: currency rate not found")

type Repository interface {
	SaveRate(name string, value float64) (model.CurrencyRate, error)
	GetRateByCurrency(currency string) (model.CurrencyRate, error)
}
