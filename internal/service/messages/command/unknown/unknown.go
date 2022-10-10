package unknown

import (
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
)

type unknownCommand struct{}

func New() messages.Command {
	return &unknownCommand{}
}

func (s *unknownCommand) Process(in model.MessageIn) (model.MessageOut, error) {
	out := model.MessageOut{
		Text: "не знаю такую команду",
	}
	return out, nil
}
