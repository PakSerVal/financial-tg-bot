package inmemory

import (
	"context"
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

func (i *inmemory) SaveRate(ctx context.Context, name string, rate int64) error {
	rateRecord := model.CurrencyRate{
		Code:  name,
		Value: rate,
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.rates[name] = rateRecord

	return nil
}

func (i *inmemory) GetRateByCurrency(ctx context.Context, currency string) (*model.CurrencyRate, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	rate, ok := i.rates[currency]
	if !ok {
		return nil, nil
	}

	return &rate, nil
}
