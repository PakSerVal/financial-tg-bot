package currency_rate

import "github.com/pkg/errors"

var ErrCurrencyRateNotFound = errors.New("repo: currency rate not found")

type Repository interface {
	SaveRate(name string, value float64) (CurrencyRate, error)
	GetRateByCurrency(currency string) (CurrencyRate, error)
}
