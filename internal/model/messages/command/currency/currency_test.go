package currency

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/dto"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/mocks"
	mockSelectedCurrency "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency/mocks"
)

func TestReportCommand_ProcessFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mockMessages.NewMockCommand(ctrl)
	repo := mockSelectedCurrency.NewMockRepository(ctrl)
	command := New(next, repo)

	gomock.InOrder(
		next.EXPECT().Process(dto.MessageIn{Text: "not_supported"}).Return(dto.MessageOut{Text: "тест"}, nil),
		repo.EXPECT().SaveSelectedCurrency("EUR", int64(1)).Return(errors.New("some error")),
		repo.EXPECT().SaveSelectedCurrency("USD", int64(2)).Return(nil),
	)

	t.Run("not supported", func(t *testing.T) {
		res, err := command.Process(dto.MessageIn{Text: "not_supported"})

		assert.NoError(t, err)
		assert.Equal(t, dto.MessageOut{Text: "тест"}, res)
	})

	t.Run("currency buttons", func(t *testing.T) {
		res, err := command.Process(dto.MessageIn{Text: "currency"})

		assert.NoError(t, err)
		assert.Equal(t, dto.MessageOut{Text: "В какой валюте вы хотите получать отчеты?", KeyBoard: &dto.KeyBoard{
			OneTime: true,
			Rows: []dto.KeyBoardRow{
				{
					Buttons: []dto.KeyBoardButton{
						{Text: Usd}, {Text: Eur}, {Text: Rub}, {Text: Cny},
					},
				},
			},
		}}, res)
	})

	t.Run("saving currency repo error", func(t *testing.T) {
		res, err := command.Process(dto.MessageIn{Text: "EUR", UserId: 1})

		assert.Error(t, err)
		assert.Equal(t, dto.MessageOut{}, res)
	})

	t.Run("saving currency repo success", func(t *testing.T) {
		res, err := command.Process(dto.MessageIn{Text: "USD", UserId: 2})

		assert.NoError(t, err)
		assert.Equal(t, dto.MessageOut{Text: "Выбранная валюта: USD"}, res)
	})
}
