package inmemory

import (
	"sync"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
)

type inmemory struct {
	rates map[string]model.CurrencyRate
	mu    *sync.Mutex
}

func New() currency_rate.Repository {
	return &inmemory{
		rates: map[string]model.CurrencyRate{},
		mu:    &sync.Mutex{},
	}
}

func (i *inmemory) SaveRate(name string, rate float64) (model.CurrencyRate, error) {
	rateRecord := model.CurrencyRate{
		Name:  name,
		Value: rate,
	}

	i.mu.Lock()
	i.rates[name] = rateRecord
	i.mu.Unlock()

	return rateRecord, nil
}

func (i *inmemory) GetRateByCurrency(currency string) (model.CurrencyRate, error) {
	i.mu.Lock()
	rate, ok := i.rates[currency]
	if !ok {
		return rate, currency_rate.ErrCurrencyRateNotFound
	}
	i.mu.Unlock()

	return rate, nil
}
