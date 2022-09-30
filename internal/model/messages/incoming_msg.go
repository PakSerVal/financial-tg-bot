package messages

import (
	"github.com/pkg/errors"
)

type Command interface {
	Process(msgText string) (string, error)
}

type MessageSender interface {
	SendMessage(text string, userID int64) error
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

func (s *Model) IncomingMessage(msg Message) error {
	msgText, err := s.commandChain.Process(msg.Text)
	if err != nil {
		return errors.Wrap(err, "process incoming message error")
	}

	return s.tgClient.SendMessage(msgText, msg.UserID)
}
