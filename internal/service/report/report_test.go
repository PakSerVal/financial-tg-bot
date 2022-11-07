package report

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	currencyRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	mockCurrencyRate "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate/mocks"
	mockSelectedCurrency "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/mocks"
	mock_cache "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/cache/mocks"
	mockSpend "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/mocks"
)

func TestReportCommand_ProcessFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	sendRepo := mockSpend.NewMockRepository(ctrl)
	currencyRateRepo := mockCurrencyRate.NewMockRepository(ctrl)
	selectedCurrencyRepo := mockSelectedCurrency.NewMockRepository(ctrl)
	cache := mock_cache.NewMockSpendRepo(ctrl)

	service := New(sendRepo, currencyRateRepo, selectedCurrencyRepo, cache)

	t.Run("send repo error", func(t *testing.T) {
		cache.EXPECT().GetByTimeSince(context.TODO(), gomock.Any(), gomock.Any()).Return(nil, nil)

		sendRepo.EXPECT().GetByTimeSince(context.TODO(), gomock.Any(), gomock.Any()).Return([]currencyRepo.Spend{}, errors.New("some error")).Times(1)
		_, err := service.MakeReport(context.TODO(), 1, time.Now(), "сегодня")

		assert.Error(t, err)
	})

	t.Run("no records", func(t *testing.T) {
		cache.EXPECT().GetByTimeSince(context.TODO(), gomock.Any(), gomock.Any()).Return(nil, nil)
		sendRepo.EXPECT().GetByTimeSince(context.TODO(), gomock.Any(), gomock.Any()).Return([]currencyRepo.Spend{}, nil).Times(1)
		cache.EXPECT().Save(context.TODO(), gomock.Any(), gomock.Any(), []currencyRepo.Spend{}).Return(nil)
		res, err := service.MakeReport(context.TODO(), 1, time.Now(), "сегодня")

		assert.NoError(t, err)
		assert.Equal(t, &currencyRepo.MessageOut{Text: "Расходов сегодня нет"}, res)
	})

	t.Run("selected currency repo error", func(t *testing.T) {
		cache.EXPECT().GetByTimeSince(context.TODO(), gomock.Any(), gomock.Any()).Return(nil, nil)
		sendRepo.EXPECT().GetByTimeSince(gomock.Any(), int64(1), gomock.Any()).Return([]currencyRepo.Spend{{Id: 1}}, nil).Times(1)
		cache.EXPECT().Save(context.TODO(), gomock.Any(), gomock.Any(), []currencyRepo.Spend{{Id: 1}}).Return(nil)
		selectedCurrencyRepo.EXPECT().GetSelectedCurrency(context.TODO(), gomock.Any()).Return(nil, errors.New("some error")).Times(1)
		_, err := service.MakeReport(context.TODO(), int64(1), time.Now(), "сегодня")

		assert.Error(t, err)
	})

	t.Run("currency rate repo error", func(t *testing.T) {
		cache.EXPECT().GetByTimeSince(context.TODO(), gomock.Any(), gomock.Any()).Return(nil, nil)
		sendRepo.EXPECT().GetByTimeSince(context.TODO(), int64(1), gomock.Any()).Return([]currencyRepo.Spend{{Id: 1}}, nil).Times(1)
		cache.EXPECT().Save(context.TODO(), gomock.Any(), gomock.Any(), []currencyRepo.Spend{{Id: 1}}).Return(nil)
		selectedCurrencyRepo.EXPECT().GetSelectedCurrency(context.TODO(), gomock.Any()).Return(&currencyRepo.SelectedCurrency{Code: "EUR"}, nil).Times(1)
		currencyRateRepo.EXPECT().GetRateByCurrency(context.TODO(), gomock.Any()).Return(nil, errors.New("some error")).Times(1)
		_, err := service.MakeReport(context.TODO(), int64(1), time.Now(), "сегодня")

		assert.Error(t, err)
	})

	t.Run("cache get error", func(t *testing.T) {
		cache.EXPECT().GetByTimeSince(context.TODO(), gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
		_, err := service.MakeReport(context.TODO(), int64(1), time.Now(), "сегодня")

		assert.Error(t, err)
	})

	t.Run("cache set error", func(t *testing.T) {
		cache.EXPECT().GetByTimeSince(context.TODO(), gomock.Any(), gomock.Any()).Return(nil, nil)
		sendRepo.EXPECT().GetByTimeSince(context.TODO(), int64(1), gomock.Any()).Return([]currencyRepo.Spend{{Id: 1}}, nil).Times(1)
		cache.EXPECT().Save(context.TODO(), gomock.Any(), gomock.Any(), []currencyRepo.Spend{{Id: 1}}).Return(errors.New("some error"))
		_, err := service.MakeReport(context.TODO(), int64(1), time.Now(), "сегодня")

		assert.Error(t, err)
	})
}

func TestReportCommand_ProcessSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	sendRepo := mockSpend.NewMockRepository(ctrl)
	selectedCurrencyRepo := mockSelectedCurrency.NewMockRepository(ctrl)
	currencyRateRepo := mockCurrencyRate.NewMockRepository(ctrl)
	cache := mock_cache.NewMockSpendRepo(ctrl)

	service := New(sendRepo, currencyRateRepo, selectedCurrencyRepo, cache)

	spends := []currencyRepo.Spend{
		{
			Id:       1,
			Price:    10000,
			Category: "Такси",
			UserId:   1,
		},
		{
			Id:       2,
			Price:    40000,
			Category: "Такси",
			UserId:   1,
		},
		{
			Id:       3,
			Price:    20000,
			Category: "Такси",
			UserId:   1,
		},
		{
			Id:       4,
			Price:    20000,
			Category: "Продукты",
			UserId:   1,
		},
		{
			Id:       4,
			Price:    90000,
			Category: "Продукты",
			UserId:   1,
		},
		{
			Id:       5,
			Price:    200000,
			Category: "Инвестиции",
			UserId:   1,
		},
	}

	cache.EXPECT().GetByTimeSince(context.TODO(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(3)
	sendRepo.EXPECT().GetByTimeSince(context.TODO(), int64(1), gomock.Any()).Return(spends, nil).Times(3)
	cache.EXPECT().Save(context.TODO(), gomock.Any(), gomock.Any(), spends).Return(nil).Times(3)

	type args struct {
		userId int64
		prefix string
	}

	tests := []struct {
		name     string
		in       args
		currency string
		rate     int64
		wanted   *currencyRepo.MessageOut
	}{
		{
			name: "today",
			in: args{
				userId: 1,
				prefix: "сегодня",
			},
			currency: "USD",
			rate:     6010,
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
			rate:     6430,
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
			rate:     876,
			wanted: &currencyRepo.MessageOut{
				Text: "Расходы в этом году:\nИнвестиции - 228.31 юан.\nПродукты - 125.57 юан.\nТакси - 79.91 юан.",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selectedCurrencyRepo.EXPECT().GetSelectedCurrency(context.TODO(), test.in.userId).Return(&currencyRepo.SelectedCurrency{
				Code:   test.currency,
				UserId: test.in.userId,
			}, nil).Times(1)

			currencyRateRepo.EXPECT().GetRateByCurrency(context.TODO(), test.currency).Return(&currencyRepo.CurrencyRate{
				Code:  test.currency,
				Value: test.rate,
			}, nil).Times(1)

			res, err := service.MakeReport(context.TODO(), test.in.userId, time.Now(), test.in.prefix)

			assert.NoError(t, err)
			assert.Equal(t, test.wanted, res)
		})
	}
}
