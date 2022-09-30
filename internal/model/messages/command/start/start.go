package start

import "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"

const menuText = "Бот для учета финансов\n\n" +
	"Добавить трату: 350 продукты\n\n" +
	"Получить отчет: \n" +
	"- за сегодня: /today\n" +
	"- за месяц: /month\n" +
	"- за год: /year\n"

type startCommand struct {
	next messages.Command
}

func New(next messages.Command) *startCommand {
	return &startCommand{
		next: next,
	}
}

func (s *startCommand) Process(msgText string) (string, error) {
	if msgText == "/start" {
		return menuText, nil
	}

	return s.next.Process(msgText)
}
