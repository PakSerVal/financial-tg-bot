package inmemory

import (
	"sync"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
)

type inmemory struct {
	rates map[string]model.CurrencyRate
	mu    sync.RWMutex
}

func New() currency_rate.Repository {
	return &inmemory{
		rates: map[string]model.CurrencyRate{},
	}
}

func (i *inmemory) SaveRate(name string, rate float64) (model.CurrencyRate, error) {

	rateRecord := model.CurrencyRate{
		Name:  name,
		Value: rate,
	}

	i.mu.RLock()
	i.rates[name] = rateRecord
	i.mu.RUnlock()

	return rateRecord, nil
}

func (i *inmemory) GetRateByCurrency(currency string) (model.CurrencyRate, error) {
	rate, ok := i.rates[currency]
	if !ok {
		return rate, currency_rate.ErrCurrencyRateNotFound
	}

	return rate, nil
}
