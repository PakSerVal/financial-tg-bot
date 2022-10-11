package report

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	currencyRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	mockCurrencyRate "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate/mocks"
	mockSelectedCurrency "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/mocks"
	mockSpend "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/mocks"
)

func TestReportCommand_ProcessFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	sendRepo := mockSpend.NewMockRepository(ctrl)
	currencyRateRepo := mockCurrencyRate.NewMockRepository(ctrl)
	selectedCurrencyRepo := mockSelectedCurrency.NewMockRepository(ctrl)

	service := New(sendRepo, currencyRateRepo, selectedCurrencyRepo)

	t.Run("send repo error", func(t *testing.T) {
		sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]currencyRepo.Spend{}, errors.New("some error")).Times(1)
		_, err := service.MakeReport(1, time.Now(), "сегодня")

		assert.Error(t, err)
	})

	t.Run("no records", func(t *testing.T) {
		sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]currencyRepo.Spend{}, nil).Times(1)
		res, err := service.MakeReport(1, time.Now(), "сегодня")

		assert.NoError(t, err)
		assert.Equal(t, &currencyRepo.MessageOut{Text: "Расходов сегодня нет"}, res)
	})

	t.Run("selected currency repo error", func(t *testing.T) {
		sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]currencyRepo.Spend{{ID: 1}}, nil).Times(1)
		selectedCurrencyRepo.EXPECT().GetSelectedCurrency(gomock.Any()).Return(currencyRepo.SelectedCurrency{}, errors.New("some error")).Times(1)
		_, err := service.MakeReport(1, time.Now(), "сегодня")

		assert.Error(t, err)
	})

	t.Run("currency rate repo error", func(t *testing.T) {
		sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]currencyRepo.Spend{{ID: 1}}, nil).Times(1)
		selectedCurrencyRepo.EXPECT().GetSelectedCurrency(gomock.Any()).Return(currencyRepo.SelectedCurrency{Currency: "EUR"}, nil).Times(1)
		currencyRateRepo.EXPECT().GetRateByCurrency(gomock.Any()).Return(currencyRepo.CurrencyRate{}, errors.New("some error")).Times(1)
		_, err := service.MakeReport(1, time.Now(), "сегодня")

		assert.Error(t, err)
	})
}

func TestReportCommand_ProcessSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	sendRepo := mockSpend.NewMockRepository(ctrl)
	selectedCurrencyRepo := mockSelectedCurrency.NewMockRepository(ctrl)
	currencyRateRepo := mockCurrencyRate.NewMockRepository(ctrl)

	service := New(sendRepo, currencyRateRepo, selectedCurrencyRepo)

	sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]currencyRepo.Spend{
		{
			ID:       1,
			Price:    100,
			Category: "Такси",
		},
		{
			ID:       2,
			Price:    400,
			Category: "Такси",
		},
		{
			ID:       3,
			Price:    200,
			Category: "Такси",
		},
		{
			ID:       4,
			Price:    200,
			Category: "Продукты",
		},
		{
			ID:       4,
			Price:    900,
			Category: "Продукты",
		},
		{
			ID:       5,
			Price:    2000,
			Category: "Инвестиции",
		},
	}, nil).Times(3)

	type args struct {
		userId int64
		prefix string
	}

	tests := []struct {
		name     string
		in       args
		currency string
		rate     float64
		wanted   *currencyRepo.MessageOut
	}{
		{
			name: "today",
			in: args{
				userId: 1,
				prefix: "сегодня",
			},
			currency: "USD",
			rate:     60.1,
			wanted: &currencyRepo.MessageOut{
				Text: "Расходы сегодня:\nИнвестиции - 33.28 дол.\nПродукты - 18.30 дол.\nТакси - 11.65 дол.",
			},
		},
		{
			name: "month",
			in: args{
				userId: 1,
				prefix: "в текущем месяце",
			},
			currency: "EUR",
			rate:     64.3,
			wanted: &currencyRepo.MessageOut{
				Text: "Расходы в текущем месяце:\nИнвестиции - 31.10 евро.\nПродукты - 17.11 евро.\nТакси - 10.89 евро.",
			},
		},
		{
			name: "year",
			in: args{
				userId: 1,
				prefix: "в этом году",
			},
			currency: "CNY",
			rate:     8.76,
			wanted: &currencyRepo.MessageOut{
				Text: "Расходы в этом году:\nИнвестиции - 228.31 юан.\nПродукты - 125.57 юан.\nТакси - 79.91 юан.",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selectedCurrencyRepo.EXPECT().GetSelectedCurrency(test.in.userId).Return(currencyRepo.SelectedCurrency{
				Currency: test.currency,
				UserId:   test.in.userId,
			}, nil).Times(1)

			currencyRateRepo.EXPECT().GetRateByCurrency(test.currency).Return(currencyRepo.CurrencyRate{
				Name:  test.currency,
				Value: test.rate,
			}, nil).Times(1)

			res, err := service.MakeReport(test.in.userId, time.Now(), test.in.prefix)

			assert.NoError(t, err)
			assert.Equal(t, test.wanted, res)
		})
	}
}
