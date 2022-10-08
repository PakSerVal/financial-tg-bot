package messages

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // <-- Вот это пакет приходится импортировать
)

type Command interface {
	Process(msgText string) (string, error)
}

type MessageSender interface {
	SendMessage(text string, userID int64) error
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

type Message struct {
	Text   string
	UserID int64
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
			err := s.processMessage(Message{
				Text:   text,
				UserID: update.Message.From.ID,
			})

			if err != nil {
				log.Println("error processing message:", err)
			}
		}
	}
}

func (s *Model) processMessage(message Message) error {
	msgText, err := s.commandChain.Process(message.Text)
	if err != nil {
		return err
	}

	return s.tgClient.SendMessage(msgText, message.UserID)
}
