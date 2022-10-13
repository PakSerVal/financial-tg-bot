package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	customErrors "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/err_msg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/utils"
)

const (
	Usd = "USD"
	Eur = "EUR"
	Rub = "RUB"
	Cny = "CNY"
)

var currencyUnitName = map[string]string{
	Usd: "дол",
	Eur: "евро",
	Rub: "руб",
	Cny: "юан",
}

type Service interface {
	MakeReport(userId int64, timeSince time.Time, timeRangePrefix string) (*model.MessageOut, error)
}

type service struct {
	spendRepo            spend.Repository
	currencyRateRepo     currency_rate.Repository
	selectedCurrencyRepo selected_currency.Repository
}

func New(spendRepo spend.Repository, currencyRateRepo currency_rate.Repository, selectedCurrencyRepo selected_currency.Repository) Service {
	return &service{
		spendRepo:            spendRepo,
		currencyRateRepo:     currencyRateRepo,
		selectedCurrencyRepo: selectedCurrencyRepo,
	}
}

func (r *service) MakeReport(userId int64, timeSince time.Time, timeRangePrefix string) (*model.MessageOut, error) {
	records, err := r.spendRepo.GetByTimeSince(timeSince)
	if err != nil {
		return nil, errors.Wrap(err, "spendRepo: get by time since")
	}

	if len(records) == 0 {
		return &model.MessageOut{
			Text: "Расходов " + timeRangePrefix + " нет",
		}, nil
	}

	cur, err := r.getSelectedCurrency(userId)
	if err != nil {
		return nil, errors.Wrap(err, "get selected currency error")
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
			return nil, errors.Wrap(err, "get currency rates error")
		}
	}

	var msgTextParts []string
	for category, sum := range groupRecords(records, rate.Value) {
		msgTextParts = append(msgTextParts, fmt.Sprintf("%s - %.2f %s.", category, sum, unitName))
	}
	sort.Strings(msgTextParts)

	return &model.MessageOut{
		Text: "Расходы " + timeRangePrefix + ":\n" + strings.Join(msgTextParts, "\n"),
	}, nil
}

func (r *service) getSelectedCurrency(userId int64) (model.SelectedCurrency, error) {
	selectedCurrency, err := r.selectedCurrencyRepo.GetSelectedCurrency(userId)

	if errors.Is(err, customErrors.CurrencyNotFound) {
		selectedCurrency.Currency = "руб"
		return selectedCurrency, nil
	}

	if err != nil {
		return selectedCurrency, err
	}

	return selectedCurrency, nil
}

func groupRecords(records []model.Spend, rate int64) map[string]float64 {
	rateF := utils.ConvertKopecksToFloat(rate)
	m := map[string]float64{}
	for _, record := range records {
		price := utils.ConvertKopecksToFloat(record.Price) / rateF
		m[record.Category] += price
	}

	return m
}
