package report

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/dto"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/mocks"
	currencyRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	mockCurrencyRate "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate/mocks"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	mockSelectedCurrency "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/mocks"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	mockSpend "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/mocks"
)

func TestReportCommand_ProcessFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mockMessages.NewMockCommand(ctrl)
	sendRepo := mockSpend.NewMockRepository(ctrl)
	selectedCurrencyRepo := mockSelectedCurrency.NewMockRepository(ctrl)
	currencyRateRepo := mockCurrencyRate.NewMockRepository(ctrl)
	command := New(next, sendRepo, selectedCurrencyRepo, currencyRateRepo)

	t.Run("not supported", func(t *testing.T) {
		next.EXPECT().Process(dto.MessageIn{Text: "not supported"}).Return(dto.MessageOut{Text: "test"}, nil)
		res, err := command.Process(dto.MessageIn{Text: "not supported"})

		assert.NoError(t, err)
		assert.Equal(t, dto.MessageOut{Text: "test"}, res)
	})

	t.Run("send repo error", func(t *testing.T) {
		sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]spend.SpendRecord{}, errors.New("some error")).Times(1)
		_, err := command.Process(dto.MessageIn{Text: "today"})

		assert.Error(t, err)
	})

	t.Run("no records", func(t *testing.T) {
		sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]spend.SpendRecord{}, nil).Times(1)
		res, err := command.Process(dto.MessageIn{Text: "today"})

		assert.NoError(t, err)
		assert.Equal(t, dto.MessageOut{Text: "Расходов сегодня нет"}, res)
	})

	t.Run("selected currency repo error", func(t *testing.T) {
		sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]spend.SpendRecord{{ID: 1}}, nil).Times(1)
		selectedCurrencyRepo.EXPECT().GetSelectedCurrency(gomock.Any()).Return(selected_currency.SelectedCurrency{}, errors.New("some error")).Times(1)
		_, err := command.Process(dto.MessageIn{Text: "today"})

		assert.Error(t, err)
	})

	t.Run("currency rate repo error", func(t *testing.T) {
		sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]spend.SpendRecord{{ID: 1}}, nil).Times(1)
		selectedCurrencyRepo.EXPECT().GetSelectedCurrency(gomock.Any()).Return(selected_currency.SelectedCurrency{Currency: "EUR"}, nil).Times(1)
		currencyRateRepo.EXPECT().GetRateByCurrency(gomock.Any()).Return(currencyRepo.CurrencyRate{}, errors.New("some error")).Times(1)
		_, err := command.Process(dto.MessageIn{Text: "today"})

		assert.Error(t, err)
	})
}

func TestReportCommand_ProcessSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mockMessages.NewMockCommand(ctrl)
	sendRepo := mockSpend.NewMockRepository(ctrl)
	selectedCurrencyRepo := mockSelectedCurrency.NewMockRepository(ctrl)
	currencyRateRepo := mockCurrencyRate.NewMockRepository(ctrl)
	command := New(next, sendRepo, selectedCurrencyRepo, currencyRateRepo)

	sendRepo.EXPECT().GetByTimeSince(gomock.Any()).Return([]spend.SpendRecord{
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

	tests := []struct {
		name     string
		in       dto.MessageIn
		currency string
		rate     float64
		wanted   dto.MessageOut
	}{
		{
			name: "today",
			in: dto.MessageIn{
				Text:   "today",
				UserId: 1,
			},
			currency: "USD",
			rate:     60.1,
			wanted: dto.MessageOut{
				Text: "Расходы сегодня:\nИнвестиции - 33.28 дол.\nПродукты - 18.30 дол.\nТакси - 11.65 дол.",
			},
		},
		{
			name: "month",
			in: dto.MessageIn{
				Text:   "month",
				UserId: 1,
			},
			currency: "EUR",
			rate:     64.3,
			wanted: dto.MessageOut{
				Text: "Расходы в текущем месяце:\nИнвестиции - 31.10 евро.\nПродукты - 17.11 евро.\nТакси - 10.89 евро.",
			},
		},
		{
			name: "year",
			in: dto.MessageIn{
				Text:   "year",
				UserId: 1,
			},
			currency: "CNY",
			rate:     8.76,
			wanted: dto.MessageOut{
				Text: "Расходы в этом году:\nИнвестиции - 228.31 юан.\nПродукты - 125.57 юан.\nТакси - 79.91 юан.",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selectedCurrencyRepo.EXPECT().GetSelectedCurrency(test.in.UserId).Return(selected_currency.SelectedCurrency{
				Currency: test.currency,
				UserId:   test.in.UserId,
			}, nil).Times(1)

			currencyRateRepo.EXPECT().GetRateByCurrency(test.currency).Return(currencyRepo.CurrencyRate{
				Name:  test.currency,
				Value: test.rate,
			}, nil).Times(1)

			res, err := command.Process(test.in)

			assert.NoError(t, err)
			assert.Equal(t, test.wanted, res)
		})
	}
}
