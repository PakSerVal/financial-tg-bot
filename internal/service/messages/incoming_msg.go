package messages

import (
	"log"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

type Command interface {
	Process(in model.MessageIn) (*model.MessageOut, error)
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

func (s *Model) ListenIncomingMessages() {
	log.Println("listening for messages")
	ch := s.tgClient.GetUpdatesChan()

	for update := range ch {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			text := update.Message.Text
			if update.Message.IsCommand() {
				text = update.Message.Command()
			}
			err := s.processMessage(text, update.Message.From.ID)

			if err != nil {
				log.Println("error processing message:", err)
			}
		}
	}
}

func (s *Model) processMessage(msgText string, userId int64) error {
	msgOut, err := s.commandChain.Process(model.MessageIn{
		Text:   msgText,
		UserId: userId,
	})
	if err != nil {
		return err
	}

	if msgOut == nil {
		return errors.New("message result must be non-empty")
	}

	return s.tgClient.SendMessage(*msgOut, userId)
}
