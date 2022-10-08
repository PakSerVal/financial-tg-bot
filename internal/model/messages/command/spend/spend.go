package spend

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

type Repository interface {
	Save(sum int64, category string) (spend.SpendRecord, error)
	GetByTimeSince(timeSince time.Time) ([]spend.SpendRecord, error)
}

type spendCommand struct {
	next messages.Command
	repo Repository
}

func New(next messages.Command, repo Repository) *spendCommand {
	return &spendCommand{
		next: next,
		repo: repo,
	}
}

func (s *spendCommand) Process(msgText string) (string, error) {
	if price, category, ok := parse(msgText); ok {
		rec, err := s.repo.Save(price, category)
		if err != nil {
			return "", errors.Wrap(err, "repo: save spend record error")
		}

		return fmt.Sprintf("Добавлена трата: %s %d руб.", rec.Category, rec.Price), nil
	}

	return s.next.Process(msgText)
}

func parse(msgText string) (int64, string, bool) {
	parts := strings.Split(msgText, " ")
	if len(parts) != 2 {
		return 0, "", false
	}

	price, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", false
	}

	return int64(price), parts[1], true
}
