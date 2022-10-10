package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/currency_rates"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command"
	currencyInmemory "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate/inmemory"
	selectedCurrencyInmemory "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/inmemory"
	spendInmemory "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/inmemory"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	currencyRepo := currencyInmemory.New()

	go func() {
		currencyRateApiClient := currency_rate.NewCurrencyRateApiClient()
		currencyRatePuller := currency_rates.NewCurrencyRatePuller(currencyRepo, currencyRateApiClient)
		currencyRatePuller.Pull(ctx)
	}()

	tgClient, err := tg.New(cfg)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	spendRepo := spendInmemory.New()
	selectedCurrencyRepo := selectedCurrencyInmemory.New()
	msgModel := messages.New(tgClient, command.MakeChain(spendRepo, currencyRepo, selectedCurrencyRepo))

	msgModel.ListenIncomingMessages()
}
