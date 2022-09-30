package spend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

type spendCommand struct {
	next messages.Command
	repo Repository
}

type Repository interface {
	Save(sum int64, category string) (spend.Record, error)
}

func New(next messages.Command, repo Repository) *spendCommand {
	return &spendCommand{
		next: next,
		repo: repo,
	}
}

func (s *spendCommand) Process(msgText string) (string, error) {
	if sum, category, ok := parse(msgText); ok {
		rec, err := s.repo.Save(sum, category)
		if err != nil {
			return "", errors.Wrap(err, "repo: save spend record error")
		}

		return fmt.Sprintf("Добавлена трата: %s %d руб.", rec.Category, rec.Sum), nil
	}

	return s.next.Process(msgText)
}

func parse(msgText string) (int64, string, bool) {
	parts := strings.Split(msgText, " ")
	if len(parts) != 2 {
		return 0, "", false
	}

	sum, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", false
	}

	return int64(sum), parts[1], true
}
