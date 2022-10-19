package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/database"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/budget"
	currencyRateRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	currencyRateDB "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate/db"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate/inmemory"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	selectedCurrencyDB "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/db"
	selectedRepoInmemory "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/inmemory"
	spendRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	spendDB "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/db"
	spendRepoInmemory "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/inmemory"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/currency_rates"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report"
)

func main() {
	var err error

	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Fatal(panicErr)
		}

		if err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	db, err := database.Connect(cfg.DbConn())
	if err != nil {
		log.Fatal("db connect failed:", err)
	}

	sqlManager := database.NewSqlManager(db)

	var currencyRepo currencyRateRepo.Repository
	var spendRepo spendRepo.Repository
	var selectedCurrencyRepo selected_currency.Repository

	if cfg.UseInmemory() {
		currencyRepo = inmemory.New()
		spendRepo = spendRepoInmemory.New()
		selectedCurrencyRepo = selectedRepoInmemory.New()
	} else {
		currencyRepo = currencyRateDB.New(db)
		spendRepo = spendDB.New(db)
		selectedCurrencyRepo = selectedCurrencyDB.New(db)
	}

	budgetRepo := budget.New(db)

	go func() {
		currencyRateApiClient := currency_rate.NewCurrencyRateApiClient()
		currencyRatePuller := currency_rates.NewCurrencyRatePuller(currencyRepo, currencyRateApiClient)
		currencyRatePuller.Pull(ctx)
	}()

	tgClient, err := tg.New(cfg)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	reportService := report.New(spendRepo, currencyRepo, selectedCurrencyRepo)
	msgModel := messages.New(tgClient, command.MakeChain(spendRepo, selectedCurrencyRepo, reportService, budgetRepo, sqlManager))

	msgModel.ListenIncomingMessages(ctx)
}
