package currency

import (
	"fmt"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
)

const (
	cmdName = "currency"
	Usd     = "USD"
	Eur     = "EUR"
	Rub     = "RUB"
	Cny     = "CNY"
)

var keyBoard = &model.KeyBoard{
	OneTime: true,
	Rows: []model.KeyBoardRow{{Buttons: []model.KeyBoardButton{
		{Text: Usd},
		{Text: Eur},
		{Text: Rub},
		{Text: Cny},
	}}},
}

type currencyCommand struct {
	next messages.Command
	repo selected_currency.Repository
}

func New(next messages.Command, repo selected_currency.Repository) messages.Command {
	return &currencyCommand{
		next: next,
		repo: repo,
	}
}

func (s *currencyCommand) Process(in model.MessageIn) (model.MessageOut, error) {
	out := model.MessageOut{}
	if in.Text == cmdName {
		out.Text = "В какой валюте вы хотите получать отчеты?"
		out.KeyBoard = keyBoard

		return out, nil
	}

	if isCurrency(in.Text) {
		err := s.repo.SaveSelectedCurrency(in.Text, in.UserId)
		if err != nil {
			return out, err
		}

		out.Text = fmt.Sprintf("Выбранная валюта: %s", in.Text)
		return out, nil
	}

	return s.next.Process(in)
}

func isCurrency(msgText string) bool {
	return msgText == Usd || msgText == Eur || msgText == Rub || msgText == Cny
}
