package unknown

import (
	"context"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
)

type unknownCommand struct{}

func New() messages.Command {
	return &unknownCommand{}
}

func (s *unknownCommand) Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error) {
	return &model.MessageOut{
		Text: "не знаю такую команду",
	}, nil
}
