package messages

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

type Command interface {
	Process(in model.MessageIn) (model.MessageOut, error)
}

type MessageSender interface {
	SendMessage(msgOut model.MessageOut, userID int64) error
	GetUpdatesChan() tgbotapi.UpdatesChannel
}

type Model struct {
	tgClient     MessageSender
	commandChain Command
}

func New(tgClient MessageSender, commandChain Command) *Model {
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

	return s.tgClient.SendMessage(msgOut, userId)
}
