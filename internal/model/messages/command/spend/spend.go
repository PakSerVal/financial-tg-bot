package spend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/dto"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

type spendCommand struct {
	next messages.Command
	repo spend.Repository
}

func New(next messages.Command, repo spend.Repository) messages.Command {
	return &spendCommand{
		next: next,
		repo: repo,
	}
}

func (s *spendCommand) Process(in dto.MessageIn) (dto.MessageOut, error) {
	out := dto.MessageOut{}
	if price, category, ok := parse(in.Text); ok {
		rec, err := s.repo.Save(price, category)
		if err != nil {
			return out, errors.Wrap(err, "repo: save spend record error")
		}

		out.Text = fmt.Sprintf("Добавлена трата: %s %.2f руб.", rec.Category, rec.Price)
		return out, nil
	}

	return s.next.Process(in)
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
