package unknown

import (
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/dto"
)

type unknownCommand struct{}

func New() messages.Command {
	return &unknownCommand{}
}

func (s *unknownCommand) Process(in dto.MessageIn) (dto.MessageOut, error) {
	out := dto.MessageOut{
		Text: "не знаю такую команду",
	}
	return out, nil
}
