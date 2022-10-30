package currency

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report"
)

const (
	cmdName = "currency"
)

var keyBoard = &model.KeyBoard{
	OneTime: true,
	Rows: []model.KeyBoardRow{{Buttons: []model.KeyBoardButton{
		{Text: report.Usd},
		{Text: report.Eur},
		{Text: report.Rub},
		{Text: report.Cny},
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

func (c *currencyCommand) Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error) {
	if in.Command == cmdName {
		return &model.MessageOut{
			Text:     "В какой валюте вы хотите получать отчеты?",
			KeyBoard: keyBoard,
		}, nil
	}

	if isCurrency(in.Command) {
		err := c.repo.SaveSelectedCurrency(ctx, in.Command, in.UserId)
		if err != nil {
			return nil, err
		}

		return &model.MessageOut{
			Text: fmt.Sprintf("Выбранная валюта: %c", in.Command),
		}, nil
	}

	return c.next.Process(ctx, in)
}

func (c *currencyCommand) Name() string {
	return cmdName
}

func isCurrency(msgText string) bool {
	return msgText == report.Usd || msgText == report.Eur || msgText == report.Rub || msgText == report.Cny
}
