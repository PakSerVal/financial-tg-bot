package queue_message

import (
	"encoding/json"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

const topic = "report"

type Sender interface {
	Send(msg model.ReportMsg) error
}

type reportSender struct {
	producer sarama.AsyncProducer
}

func NewSender(producer sarama.AsyncProducer) Sender {
	return &reportSender{producer: producer}
}

func (r *reportSender) Send(msg model.ReportMsg) error {
	encoded, err := json.Marshal(msg)
	if err != nil {
		return errors.WithStack(err)
	}

	kafkaMsg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strconv.Itoa(int(msg.UserId))),
		Value: sarama.StringEncoder(encoded),
	}
	r.producer.Input() <- &kafkaMsg

	return nil
}
