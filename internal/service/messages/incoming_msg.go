package messages

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/metrics"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"go.uber.org/zap"
)

type Command interface {
	Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error)
	Name() string
}

type Model struct {
	tgClient     tg.Client
	commandChain Command
}

func New(tgClient tg.Client, commandChain Command) *Model {
	return &Model{
		tgClient:     tgClient,
		commandChain: commandChain,
	}
}

func (s *Model) ListenIncomingMessages(ctx context.Context) {
	logger.Info("listening for messages")
	ch := s.tgClient.GetUpdatesChan()

	for update := range ch {
		if update.Message == nil {
			continue
		}

		span, ctx := opentracing.StartSpanFromContext(
			ctx,
			"incoming message",
		)

		logger.Info("incoming message",
			zap.String("user_name", update.Message.From.UserName),
			zap.String("message test", update.Message.Text),
		)

		command := update.Message.Text
		if update.Message.IsCommand() {
			command = update.Message.Command()
		}

		err := s.processMessage(ctx, command, update.Message.CommandArguments(), update.Message.From.ID)

		if err != nil {
			logger.Error("error processing message:", zap.Error(err))
			metrics.IncomingMessageTotal("failed")
		} else {
			metrics.IncomingMessageTotal("success")
		}

		span.Finish()
	}
}

func (s *Model) processMessage(ctx context.Context, command string, arguments string, userId int64) error {
	msgOut, err := s.commandChain.Process(ctx, model.MessageIn{
		Command:   command,
		Arguments: arguments,
		UserId:    userId,
	})
	if err != nil {
		return err
	}

	if msgOut != nil {
		return s.tgClient.SendMessage(*msgOut, userId)
	}

	return nil
}
