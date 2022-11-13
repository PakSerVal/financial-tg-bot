package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"go.uber.org/zap"
)

type Handler interface {
	HandleMessage(ctx context.Context, msg model.ReportMsg)
}

type reader struct {
	cfg     *sarama.Config
	cg      sarama.ConsumerGroup
	topics  []string
	handler Handler
	cancel  context.CancelFunc
}

func NewReader(
	kafkaConfig config.KafkaConfig,
	handler Handler,
) (*reader, error) {
	var err error

	cfg := sarama.NewConfig()

	ver, err := sarama.ParseKafkaVersion(kafkaConfig.Version)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cfg.Version = ver

	cfg.Consumer.Offsets.Initial = kafkaConfig.Consumer.Offset

	switch kafkaConfig.Consumer.Assignor {
	case "sticky":
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "round-robin":
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		logger.Fatal("Unrecognized consumer group partition assignor", zap.String("assignor", kafkaConfig.Consumer.Assignor))
	}

	r := &reader{
		topics:  strings.Split(kafkaConfig.Consumer.Topics, ","),
		handler: handler,
		cfg:     cfg,
	}

	r.cg, err = sarama.NewConsumerGroup(strings.Split(kafkaConfig.BrokerList, ","), kafkaConfig.Consumer.GroupId, r.cfg)
	if err != nil {
		return nil, fmt.Errorf("create consumer group: %w", err)
	}

	return r, nil
}

func (m *reader) Start(ctx context.Context) {
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := m.cg.Consume(
			ctx,
			m.topics,
			&consumer{handler: m.handler},
		); err != nil {
			logger.Fatal("Error from consumer", zap.Error(err))
		}

		if ctx.Err() != nil {
			return
		}
	}
}

// consumer implements sarama.ConsumerGroupHandler
type consumer struct {
	handler Handler
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	//
	// (multiple goroutines actually)
	ctx := session.Context()

	for msg := range claim.Messages() {
		var reportMsg model.ReportMsg
		err := json.Unmarshal(msg.Value, &reportMsg)
		if err != nil {
			logger.Error("kafka message unmarshal error", zap.Error(err))
		}

		c.handler.HandleMessage(ctx, reportMsg)

		session.MarkMessage(msg, "")
	}

	return nil
}
