package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/database"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/budget"
	currencyRateDB "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate/db"
	selectedCurrencyDB "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/db"
	spendDB "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/db"
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
	currencyRepo := currencyRateDB.New(db)

	go func() {
		currencyRateApiClient := currency_rate.NewCurrencyRateApiClient()
		currencyRatePuller := currency_rates.NewCurrencyRatePuller(currencyRepo, currencyRateApiClient)
		currencyRatePuller.Pull(ctx)
	}()

	tgClient, err := tg.New(cfg)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	spendRepo := spendDB.New(db)
	budgetRepo := budget.New(db)
	selectedCurrencyRepo := selectedCurrencyDB.New(db)
	reportService := report.New(spendRepo, currencyRepo, selectedCurrencyRepo)
	msgModel := messages.New(tgClient, command.MakeChain(spendRepo, selectedCurrencyRepo, reportService, budgetRepo, sqlManager))

	msgModel.ListenIncomingMessages(ctx)
}
