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

func (i *inmemory) SaveRate(name string, rate int64) (model.CurrencyRate, error) {
	rateRecord := model.CurrencyRate{
		Name:  name,
		Value: rate,
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.rates[name] = rateRecord

	return rateRecord, nil
}

func (i *inmemory) GetRateByCurrency(currency string) (model.CurrencyRate, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	rate, ok := i.rates[currency]
	if !ok {
		return rate, currency_rate.ErrCurrencyRateNotFound
	}

	return rate, nil
}
