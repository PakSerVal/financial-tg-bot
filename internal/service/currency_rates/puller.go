package currency_rates

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	currencyRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/utils"
	"go.uber.org/zap"
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
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"pulling rates",
	)
	defer span.Finish()

	logger.Info("pulling rates...")
	rates, err := c.apiClient.GetCurrencyRates()
	if err != nil {
		logger.Error("getting rates error", zap.Error(err))
		return
	}

	for _, rate := range rates {
		err = c.currencyRepo.SaveRate(ctx, rate.Name, utils.ConvertFloatToKopecks(rate.Rate))
		if err != nil {
			logger.Error("saving rate to db error:", zap.Error(err))
			return
		}
	}
}
