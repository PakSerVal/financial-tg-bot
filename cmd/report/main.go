package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/database"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/grpc/interceptors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/kafka"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/redis"
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
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report/queue_message"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/tracing"
	api "gitlab.ozon.dev/paksergey94/telegram-bot/pkg"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.GrpcHost(), cfg.GrpcPort()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.ClientMetricsInterceptor),
	)
	if err != nil {
		logger.Fatal("can not connect grpc", zap.Error(err))
	}

	redisClient := redis.NewClient(cfg.RedisConfig())
	spendCache := redisRepo.New(redisClient)

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

	reportClient := api.NewReportClient(conn)
	reportService := report.New(spendRepo, currencyRepo, selectedCurrencyRepo, spendCache)

	msgHandler := queue_message.NewHandler(reportClient, reportService)

	consumer, err := kafka.NewReader(cfg.KafkaConfig(), msgHandler)
	if err != nil {
		logger.Fatal("can not connect kafka consumer", zap.Error(err))
	}

	consumer.Start(ctx)
}
