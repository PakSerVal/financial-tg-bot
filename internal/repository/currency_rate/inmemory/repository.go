package inmemory

import (
	"sync"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
)

type inmemory struct {
	rates map[string]currency_rate.CurrencyRate
	mu    *sync.Mutex
}

func New() currency_rate.Repository {
	return &inmemory{
		rates: map[string]currency_rate.CurrencyRate{},
		mu:    &sync.Mutex{},
	}
}

func (i *inmemory) SaveRate(name string, rate float64) (currency_rate.CurrencyRate, error) {
	rateRecord := currency_rate.CurrencyRate{
		Name:  name,
		Value: rate,
	}

	i.mu.Lock()
	i.rates[name] = rateRecord
	i.mu.Unlock()

	return rateRecord, nil
}

func (i *inmemory) GetRateByCurrency(currency string) (currency_rate.CurrencyRate, error) {
	i.mu.Lock()
	rate, ok := i.rates[currency]
	if !ok {
		return rate, currency_rate.ErrCurrencyRateNotFound
	}
	i.mu.Unlock()

	return rate, nil
}
