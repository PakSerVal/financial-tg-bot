package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	appConfig "gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	"go.uber.org/zap"
)

func InitTracing(appCfg *appConfig.Config) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.InitGlobalTracer(appCfg.ServiceName())
	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}
}
