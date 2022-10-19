package currency_rates

import (
	"context"
	"log"
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/currency_rate"
	currencyRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/utils"
)

const pullingInterval = time.Hour * 12

type CurrencyRatePuller struct {
	currencyRepo currencyRepo.Repository
	apiClient    currency_rate.CurrencyRateApiClient
}

func NewCurrencyRatePuller(
	currencyRepo currencyRepo.Repository,
	apiClient currency_rate.CurrencyRateApiClient,
) *CurrencyRatePuller {
	return &CurrencyRatePuller{
		currencyRepo: currencyRepo,
		apiClient:    apiClient,
	}
}

func (c *CurrencyRatePuller) Pull(ctx context.Context) {
	c.updateRates(ctx)

	ticker := time.NewTicker(pullingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			c.updateRates(ctx)
		}
	}
}

func (c *CurrencyRatePuller) updateRates(ctx context.Context) {
	log.Println("pulling rates...")
	rates, err := c.apiClient.GetCurrencyRates()
	if err != nil {
		log.Println("getting rates error:", err)
		return
	}

	for _, rate := range rates {
		err = c.currencyRepo.SaveRate(ctx, rate.Name, utils.ConvertFloatToKopecks(rate.Rate))
		if err != nil {
			log.Println("saving rate to db error:", err)
			return
		}
	}
}
