package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/inmemory"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/command/currency"
)

const (
	commandToday = "today"
	commandMonth = "month"
	commandYear  = "year"
)

var currencyUnitName = map[string]string{
	currency.Usd: "дол",
	currency.Rub: "руб",
	currency.Cny: "юан",
	currency.Eur: "евро",
}

type reportCommand struct {
	next                 messages.Command
	repo                 spend.Repository
	selectedCurrencyRepo selected_currency.Repository
	currencyRateRepo     currency_rate.Repository
}

func New(
	next messages.Command,
	repo spend.Repository,
	selectedCurrencyRepo selected_currency.Repository,
	currencyRateRepo currency_rate.Repository,
) messages.Command {
	return &reportCommand{
		next:                 next,
		repo:                 repo,
		selectedCurrencyRepo: selectedCurrencyRepo,
		currencyRateRepo:     currencyRateRepo,
	}
}

func (r *reportCommand) Process(in model.MessageIn) (model.MessageOut, error) {
	now := time.Now()
	switch in.Text {
	case commandToday:
		return r.makeReport(
			in.UserId,
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			"сегодня",
		)
	case commandMonth:
		return r.makeReport(
			in.UserId,
			time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
			"в текущем месяце",
		)
	case commandYear:
		return r.makeReport(
			in.UserId,
			time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()),
			"в этом году",
		)
	}

	return r.next.Process(in)
}

func (r *reportCommand) makeReport(userId int64, timeSince time.Time, timeRangePrefix string) (model.MessageOut, error) {
	out := model.MessageOut{}
	records, err := r.repo.GetByTimeSince(timeSince)
	if err != nil {
		return out, errors.Wrap(err, "repo: get by time since")
	}

	if len(records) == 0 {
		out.Text = "Расходов " + timeRangePrefix + " нет"
		return out, nil
	}

	cur, err := r.getSelectedCurrency(userId)
	if err != nil {
		return out, errors.Wrap(err, "get selected currency error")
	}

	unitName, ok := currencyUnitName[cur.Currency]
	if !ok {
		unitName = "руб"
	}

	rate, err := r.currencyRateRepo.GetRateByCurrency(cur.Currency)
	if errors.Is(err, currency_rate.ErrCurrencyRateNotFound) {
		rate.Value = 1
	} else {
		if err != nil {
			return out, errors.Wrap(err, "get currency rates error")
		}
	}

	var msgTextParts []string
	for category, sum := range groupRecords(records, rate.Value) {
		msgTextParts = append(msgTextParts, fmt.Sprintf("%s - %.2f %s.", category, sum, unitName))
	}
	sort.Strings(msgTextParts)

	out.Text = "Расходы " + timeRangePrefix + ":\n" + strings.Join(msgTextParts, "\n")
	return out, nil
}

func (r *reportCommand) getSelectedCurrency(userId int64) (model.SelectedCurrency, error) {
	selectedCurrency, err := r.selectedCurrencyRepo.GetSelectedCurrency(userId)

	if errors.Is(err, inmemory.CurrencyNotFound) {
		selectedCurrency.Currency = "руб"
		return selectedCurrency, nil
	}

	if err != nil {
		return selectedCurrency, err
	}

	return selectedCurrency, nil
}

func groupRecords(records []model.Spend, rate float64) map[string]float64 {
	m := map[string]float64{}
	for _, record := range records {
		price := record.Price / rate
		m[record.Category] += price
	}

	return m
}
