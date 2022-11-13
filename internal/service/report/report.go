package report

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/cache"
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
	MakeReport(ctx context.Context, userId int64, timeSince time.Time, timeRangePrefix string) (string, error)
}

type service struct {
	spendRepo            spend.Repository
	currencyRateRepo     currency_rate.Repository
	selectedCurrencyRepo selected_currency.Repository
	spendCache           cache.SpendRepo
}

func New(
	spendRepo spend.Repository,
	currencyRateRepo currency_rate.Repository,
	selectedCurrencyRepo selected_currency.Repository,
	spendCache cache.SpendRepo,
) Service {
	return &service{
		spendRepo:            spendRepo,
		currencyRateRepo:     currencyRateRepo,
		selectedCurrencyRepo: selectedCurrencyRepo,
		spendCache:           spendCache,
	}
}

func (r *service) MakeReport(ctx context.Context, userId int64, timeSince time.Time, timeRangePrefix string) (string, error) {
	records, err := r.spendCache.GetByTimeSince(ctx, userId, timeSince)
	if err != nil {
		return "", errors.Wrap(err, "cache: get spends by time since")
	}

	if records == nil {
		records, err = r.spendRepo.GetByTimeSince(ctx, userId, timeSince)
		if err != nil {
			return "", errors.Wrap(err, "spendRepo: get by time since")
		}

		err = r.spendCache.Save(ctx, userId, timeSince, records)
		if err != nil {
			return "", errors.Wrap(err, "cache: saving spends error")
		}
	}

	if len(records) == 0 {
		return "Расходов " + timeRangePrefix + " нет", nil
	}

	cur, err := r.getSelectedCurrency(ctx, userId)
	if err != nil {
		return "", errors.Wrap(err, "get selected currency error")
	}

	unitName, ok := currencyUnitName[cur]
	if !ok {
		unitName = "руб"
	}

	rate, err := r.currencyRateRepo.GetRateByCurrency(ctx, cur)
	if err != nil {
		return "", errors.Wrap(err, "get currency rates error")
	}

	rateValue := int64(100)
	if rate != nil {
		rateValue = rate.Value
	}

	var msgTextParts []string
	for category, sum := range groupRecords(records, rateValue) {
		msgTextParts = append(msgTextParts, fmt.Sprintf("%s - %.2f %s.", category, sum, unitName))
	}
	sort.Strings(msgTextParts)

	return "Расходы " + timeRangePrefix + ":\n" + strings.Join(msgTextParts, "\n"), nil
}

func (r *service) getSelectedCurrency(ctx context.Context, userId int64) (string, error) {
	selectedCurrency, err := r.selectedCurrencyRepo.GetSelectedCurrency(ctx, userId)

	if errors.Is(err, sql.ErrNoRows) {
		return "руб", nil
	}

	if err != nil {
		return "", err
	}

	if selectedCurrency == nil {
		return "руб", nil
	}

	return selectedCurrency.Code, nil
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
