package start

import (
	"context"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
)

const cmdName = "start"

const menuText = "Бот для учета финансов\n\n" +
	"Добавить трату: 350 продукты\n\n" +
	"Изменить валюту: /currency\n\n" +
	"Установить месячный бюджет: /budget 30000\n\n" +
	"Получить отчет: \n" +
	"- за сегодня: /today\n" +
	"- за месяц: /month\n" +
	"- за год: /year\n"

type startCommand struct {
	next messages.Command
}

func New(next messages.Command) messages.Command {
	return &startCommand{
		next: next,
	}
}

func (s *startCommand) Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error) {
	if in.Command == cmdName {
		return &model.MessageOut{
			Text: menuText,
		}, nil
	}

	return s.next.Process(ctx, in)
}
