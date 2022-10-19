package messages

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

type Command interface {
	Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error)
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
	log.Println("listening for messages")
	ch := s.tgClient.GetUpdatesChan()

	for update := range ch {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			coommand := update.Message.Text
			if update.Message.IsCommand() {
				coommand = update.Message.Command()
			}

			err := s.processMessage(ctx, coommand, update.Message.CommandArguments(), update.Message.From.ID)

			if err != nil {
				log.Println("error processing message:", err)
			}
		}
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

	if msgOut == nil {
		return errors.New("message result must be non-empty")
	}

	return s.tgClient.SendMessage(*msgOut, userId)
}
