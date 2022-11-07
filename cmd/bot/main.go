package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/database"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/redis"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/budget"
	currencyRateRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	currencyRateDB "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate/db"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate/inmemory"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	selectedCurrencyDB "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/db"
	selectedRepoInmemory "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/inmemory"
	spendRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	redisRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/cache"
	spendDB "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/db"
	spendRepoInmemory "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/inmemory"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/currency_rates"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/tracing"
	"go.uber.org/zap"
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

	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	logger.InitLogger(cfg)
	tracing.InitTracing(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.Connect(cfg.DbConn())
	if err != nil {
		logger.Fatal("db connect failed:", zap.Error(err))
	}

	redisClient := redis.NewClient(cfg.RedisConfig())

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

	spendCache := redisRepo.New(redisClient)

	go func() {
		http.Handle("/metrics", promhttp.Handler())

		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort()), nil)
		if err != nil {
			logger.Fatal("error starting http server", zap.Error(err))
		}
	}()

	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatal("db connect failed:", zap.Error(err))
	}

	reportService := report.New(spendRepo, currencyRepo, selectedCurrencyRepo, spendCache)
	msgModel := messages.New(tgClient, command.MakeChain(spendRepo, selectedCurrencyRepo, reportService, budgetRepo, sqlManager, spendCache))

	msgModel.ListenIncomingMessages(ctx)
}
