package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/database"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/generated/api"
	grpcServer "gitlab.ozon.dev/paksergey94/telegram-bot/internal/grpc"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/grpc/interceptors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/kafka"
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
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report/queue_message"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatal("db connect failed:", zap.Error(err))
	}

	startGrpcGateway(ctx, cfg, tgClient)

	kafkaProducer, err := kafka.AsyncProducer(cfg.KafkaConfig())
	if err != nil {
		logger.Fatal("kafka producer connect failed:", zap.Error(err))
	}
	defer kafkaProducer.AsyncClose()

	msgModel := messages.New(tgClient, command.MakeChain(
		spendRepo,
		selectedCurrencyRepo,
		queue_message.NewSender(kafkaProducer),
		budgetRepo,
		sqlManager,
		spendCache,
	))

	msgModel.ListenIncomingMessages(ctx)
}

func startGrpcGateway(ctx context.Context, cfg *config.Config, tgClient tg.Client) {
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		logger.Fatal("failed to listen grpc", zap.Error(err))
	}

	server := grpcServer.New(tgClient)

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.ServerMetricsInterceptor))
	api.RegisterReportServer(s, server)

	rmux := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.Handle("/metrics", promhttp.Handler())

	err = api.RegisterReportHandlerServer(ctx, rmux, server)
	if err != nil {
		logger.Fatal("register report handler error", zap.Error(err))
	}

	httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.HttpPort()))
	if err != nil {
		log.Fatalf("failed to listen http: %v", err)
	}

	go func() {
		logger.Info("Serving grpc address", zap.Int64("port", cfg.GrpcPort()))
		reflection.Register(s)
		if err := s.Serve(grpcListener); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	go func() {
		logger.Info("Serving http address", zap.Int64("port", cfg.HttpPort()))
		err = http.Serve(httpListener, mux)
		if err != nil {
			logger.Error("http server error", zap.Error(err))
		}
	}()
}
