package kafka

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	"go.uber.org/zap"
)

func AsyncProducer(kafkaCfg config.KafkaConfig) (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(kafkaCfg.Version)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cfg.Version = version

	// So we can know the partition and offset of messages.
	cfg.Producer.Return.Successes = kafkaCfg.Producer.ReturnSuccesses

	producer, err := sarama.NewAsyncProducer(strings.Split(kafkaCfg.BrokerList, ","), cfg)
	if err != nil {
		return nil, fmt.Errorf("starting Sarama producer: %w", err)
	}

	// We will log to STDOUT if we're not able to produce messages.
	go func() {
		for err := range producer.Errors() {
			logger.Error("Failed to write queue_message", zap.Error(err))
		}
	}()

	if cfg.Producer.Return.Successes {
		go func() {
			successMsg := <-producer.Successes()
			logger.Info("Successful to write message", zap.String("topic", successMsg.Topic), zap.Int64("offset", successMsg.Offset))
		}()
	}

	return producer, nil
}
