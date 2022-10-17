package currency

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	mockSelectedCurrency "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/mocks"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/mocks"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report"
)

func TestReportCommand_ProcessFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mockMessages.NewMockCommand(ctrl)
	repo := mockSelectedCurrency.NewMockRepository(ctrl)
	command := New(next, repo)

	gomock.InOrder(
		next.EXPECT().Process(context.TODO(), model.MessageIn{Command: "not_supported"}).Return(&model.MessageOut{Text: "тест"}, nil),
		repo.EXPECT().SaveSelectedCurrency(context.TODO(), "EUR", int64(1)).Return(errors.New("some error")),
		repo.EXPECT().SaveSelectedCurrency(context.TODO(), "USD", int64(2)).Return(nil),
	)

	t.Run("not supported", func(t *testing.T) {
		res, err := command.Process(context.TODO(), model.MessageIn{Command: "not_supported"})

		assert.NoError(t, err)
		assert.Equal(t, &model.MessageOut{Text: "тест"}, res)
	})

	t.Run("currency buttons", func(t *testing.T) {
		res, err := command.Process(context.TODO(), model.MessageIn{Command: "currency"})

		assert.NoError(t, err)
		assert.Equal(t, &model.MessageOut{Text: "В какой валюте вы хотите получать отчеты?", KeyBoard: &model.KeyBoard{
			OneTime: true,
			Rows: []model.KeyBoardRow{
				{
					Buttons: []model.KeyBoardButton{
						{Text: report.Usd}, {Text: report.Eur}, {Text: report.Rub}, {Text: report.Cny},
					},
				},
			},
		}}, res)
	})

	t.Run("saving currency repo error", func(t *testing.T) {
		res, err := command.Process(context.TODO(), model.MessageIn{Command: "EUR", UserId: 1})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("saving currency repo success", func(t *testing.T) {
		res, err := command.Process(context.TODO(), model.MessageIn{Command: "USD", UserId: 2})

		assert.NoError(t, err)
		assert.Equal(t, &model.MessageOut{Text: "Выбранная валюта: USD"}, res)
	})
}
