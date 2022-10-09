package start

import (
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/dto"
)

const cmdName = "start"

const menuText = "Бот для учета финансов\n\n" +
	"Добавить трату: 350 продукты\n\n" +
	"Изменить валюту: /currency\n\n" +
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

func (s *startCommand) Process(in dto.MessageIn) (dto.MessageOut, error) {
	out := dto.MessageOut{}
	if in.Text == cmdName {
		out.Text = menuText
		return out, nil
	}

	return s.next.Process(in)
}
