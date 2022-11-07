package spend

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/database"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/budget"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/cache"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/utils"
)

type spendCommand struct {
	next       messages.Command
	spendRepo  spend.Repository
	sqlManager database.SqlManager
	budgetRepo budget.Repository
	cache      cache.SpendRepo
}

func New(next messages.Command, spendRepo spend.Repository, budgetRepo budget.Repository, manager database.SqlManager, cache cache.SpendRepo) messages.Command {
	return &spendCommand{
		next:       next,
		spendRepo:  spendRepo,
		budgetRepo: budgetRepo,
		sqlManager: manager,
		cache:      cache,
	}
}

func (s *spendCommand) Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error) {
	if price, category, ok := parse(in.Command); ok {
		var msgOut *model.MessageOut

		err := s.sqlManager.InTransaction(ctx, func(tx *sql.Tx, ctx context.Context) (bool, error) {
			err := s.spendRepo.SaveTx(tx, ctx, utils.ConvertFloatToKopecks(price), category, in.UserId)
			if err != nil {
				return false, errors.Wrap(err, "spendRepo: save spend record error")
			}

			spends, err := s.spendRepo.GetByTimeSinceTx(tx, ctx, in.UserId, utils.BeginOfMonth(time.Now()))
			if err != nil {
				return false, errors.Wrap(err, "spendRepo: getting spends error")
			}

			budg, err := s.budgetRepo.GetBudgetTx(tx, ctx, in.UserId)
			if err != nil {
				return false, err
			}

			if isBudgetLimitExceeded(spends, budg) {
				msgOut = &model.MessageOut{
					Text: fmt.Sprintf(
						"Трата не была добавлена, так как превышен лимит за текущий месяц в %.2f руб",
						utils.ConvertKopecksToFloat(budg.Value),
					),
				}

				return false, nil
			}

			msgOut = &model.MessageOut{
				Text: fmt.Sprintf("Добавлена трата: %s %.2f руб.", category, price),
			}

			return true, nil
		})
		if err != nil {
			return nil, err
		}

		err = s.cache.DeleteForUser(ctx, in.UserId)
		if err != nil {
			return nil, errors.Wrap(err, "cache: deleting spends for user error")
		}

		return msgOut, nil
	}

	return s.next.Process(ctx, in)
}

func (s *spendCommand) Name() string {
	return "spend"
}

func parse(msgText string) (float64, string, bool) {
	parts := strings.Split(msgText, " ")
	if len(parts) != 2 {
		return 0, "", false
	}

	price, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, "", false
	}

	return price, parts[1], true
}

func isBudgetLimitExceeded(spends []model.Spend, budget *model.Budget) bool {
	if budget == nil {
		return false
	}

	var sum int64
	for _, s := range spends {
		sum += s.Price
	}

	return sum > budget.Value
}
