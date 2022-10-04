package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	spendRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

const (
	CommandToday = "today"
	CommandMonth = "month"
	CommandYear  = "year"
)

type reportCommand struct {
	next messages.Command
	repo spendRepo.Repository
}

func New(next messages.Command, repo spendRepo.Repository) *reportCommand {
	return &reportCommand{
		next: next,
		repo: repo,
	}
}

func (s *reportCommand) Process(msgText string) (string, error) {
	now := time.Now()
	switch msgText {
	case CommandToday:
		return s.makeReport(
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			"сегодня",
		)
	case CommandMonth:
		return s.makeReport(
			time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
			"в текущем месяце",
		)
	case CommandYear:
		return s.makeReport(
			time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()),
			"в этом году",
		)
	}

	return s.next.Process(msgText)
}

func (s *reportCommand) makeReport(timeSince time.Time, timeRangePrefix string) (string, error) {
	records, err := s.repo.GetByTimeSince(timeSince)
	if err != nil {
		return "", errors.Wrap(err, "repo: get by time since")
	}

	if len(records) == 0 {
		return "Расходов " + timeRangePrefix + " нет", err
	}

	var msgTextParts []string
	for category, sum := range groupRecords(records) {
		msgTextParts = append(msgTextParts, fmt.Sprintf("%s - %d руб.", category, sum))
	}
	sort.Strings(msgTextParts)

	return "Расходы " + timeRangePrefix + ":\n" + strings.Join(msgTextParts, "\n"), nil
}

func groupRecords(records []spend.SpendRecord) map[string]int64 {
	m := map[string]int64{}
	for _, record := range records {
		m[record.Category] += record.Price
	}

	return m
}
