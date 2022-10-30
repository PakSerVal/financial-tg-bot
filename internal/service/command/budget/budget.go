package budget

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/budget"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/utils"
)

const (
	cmdName = "budget"
)

type budgetCommand struct {
	next messages.Command
	repo budget.Repository
}

func New(next messages.Command, repo budget.Repository) messages.Command {
	return &budgetCommand{
		next: next,
		repo: repo,
	}
}

func (b *budgetCommand) Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error) {
	if in.Command == cmdName {
		limit, err := strconv.ParseFloat(in.Arguments, 64)
		if err != nil {
			return b.next.Process(ctx, in)
		}

		err = b.repo.SaveBudget(ctx, in.UserId, utils.ConvertFloatToKopecks(limit))
		if err != nil {
			return nil, errors.Wrap(err, "saving budget error")
		}

		return &model.MessageOut{
			Text: fmt.Sprintf("Бюджет в %.2f руб установлен", limit),
		}, nil
	}

	return b.next.Process(ctx, in)
}

func (b *budgetCommand) Name() string {
	return cmdName
}
